package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/mailer"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/metrics"
	"github.com/MamangRust/monolith-ecommerce-pkg/emailretry"
	"github.com/MamangRust/monolith-ecommerce-pkg/event"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/outbox"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// retryPublisher is implemented by *kafka.Kafka and lets the handlers publish
// retry/DLQ messages with metadata headers (Phase 4) while propagating the
// active trace context (Phase 5).
type retryPublisher interface {
	SendMessageWithHeaders(ctx context.Context, topic, key string, value []byte, headers []sarama.RecordHeader) error
}

type emailHandler struct {
	ctx             context.Context
	trace           trace.Tracer
	logger          logger.LoggerInterface
	Mailer          *mailer.Mailer
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	inbox           outbox.ConsumerInbox
	consumerName    string
	producer        retryPublisher
	retryBackoff    time.Duration
}

func NewEmailHandler(ctx context.Context, logger logger.LoggerInterface, mailer *mailer.Mailer) *emailHandler {
	return newEmailHandler(ctx, logger, mailer, nil, "", nil, emailretry.DefaultBackoff)
}

// NewEmailHandlerWithInbox enables durable PostgreSQL-backed deduplication
// and retry-topic offloading for transient SMTP failures.
func NewEmailHandlerWithInbox(ctx context.Context, logger logger.LoggerInterface, mailer *mailer.Mailer, inbox outbox.ConsumerInbox, consumerName string, producer retryPublisher, retryBackoff time.Duration) *emailHandler {
	return newEmailHandler(ctx, logger, mailer, inbox, consumerName, producer, retryBackoff)
}

func newEmailHandler(ctx context.Context, logger logger.LoggerInterface, mailer *mailer.Mailer, inbox outbox.ConsumerInbox, consumerName string, producer retryPublisher, retryBackoff time.Duration) *emailHandler {
	meter := otel.Meter("email-service")

	requestCounter, err := meter.Int64Counter(
		"email_service_requests_total",
		metric.WithDescription("Total number of requests to the EmailService"),
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(err)
	}

	requestDuration, err := meter.Float64Histogram(
		"email_service_request_duration_seconds",
		metric.WithDescription("Histogram of request durations for the EmailService"),
		metric.WithUnit("s"),
	)
	if err != nil {
		panic(err)
	}

	return &emailHandler{
		ctx:             ctx,
		logger:          logger,
		Mailer:          mailer,
		trace:           otel.Tracer("email-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
		inbox:           inbox,
		consumerName:    consumerName,
		producer:        producer,
		retryBackoff:    retryBackoff,
	}
}

func (h *emailHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *emailHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *emailHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	start := time.Now()
	status := "success"

	defer func() {
		h.recordMetrics("ConsumeClaim", status, start)
	}()

	
	for msg := range claim.Messages() {
		ctx := otel.GetTextMapPropagator().Extract(h.ctx, kafkaHeaderCarrier(msg.Headers))
		ctx, span := h.trace.Start(ctx, "consume:"+msg.Topic)
		span.SetAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", msg.Topic),
			attribute.Int64("messaging.kafka.partition", int64(msg.Partition)),
			attribute.Int64("messaging.kafka.offset", msg.Offset),
		)

		var payload event.EmailEnvelope
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UNMARSHAL_MESSAGE")

			h.logger.Error("Failed to unmarshal message", zap.Error(err))

			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to unmarshal message")
			status = "failed_unmarshal_message"

			metrics.EmailInvalid.Add(ctx, 1)
			span.End()
			sess.MarkMessage(msg, "")
			sess.Commit()
			continue
		}

		if !payload.IsValid() {
			h.logger.Error("Invalid email envelope: event_id, schema_version=1, event_type, email, subject and body are required")
			status = "failed_invalid_message"

			metrics.EmailInvalid.Add(ctx, 1)
			span.End()
			sess.MarkMessage(msg, "")
			sess.Commit()
			continue
		}

		// Dedupe on the envelope's event_id (unique per event) rather than the
		// Kafka message key (a business entity ID): two distinct events for the
		// same entity must both be delivered.
		eventKey := fmt.Sprintf("%s:%s", msg.Topic, payload.EventID)
		var (
			reserved, processed bool
			reservationVersion  int64
			err                 error
		)
		if h.inbox != nil {
			reserved, processed, reservationVersion, err = h.inbox.Reserve(ctx, h.consumerName, eventKey, msg.Topic, msg.Partition, msg.Offset)
			if err != nil {
				status = "failed_inbox_reserve"
				span.End()
				return err
			}
			if processed {
				// The side effect was already completed by this consumer.
				span.End()
				sess.MarkMessage(msg, "")
				sess.Commit()
				continue
			}
			if !reserved {
				status = "inbox_lease_busy"
				span.End()
				return fmt.Errorf("consumer inbox lease is active for event %s", eventKey)
			}
		}

		err = h.Mailer.Send(payload.Email, payload.Subject, payload.Body)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_SEND_EMAIL")

			h.logger.Error("Failed to send email; offloading to retry topic", zap.Error(err))
			if h.inbox != nil {
				if releaseErr := h.inbox.Release(ctx, h.consumerName, eventKey, reservationVersion, err); releaseErr != nil {
					h.logger.Error("failed to release consumer inbox lease", zap.Error(releaseErr), zap.String("event_key", eventKey))
				}
			}
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to send email")
			status = "failed_send_email"

			metrics.EmailFailed.Add(ctx, 1)

			if h.producer == nil {
				span.End()
				return err
			}
			meta := emailretry.RetryMeta{
				SrcTopic: msg.Topic, SrcPartition: msg.Partition, SrcOffset: msg.Offset,
				Attempt: 1, RetryAt: emailretry.NextRetryAt(1, h.retryBackoff), Reason: err.Error(),
			}
			if pubErr := h.producer.SendMessageWithHeaders(ctx, emailretry.RetryTopic, payload.EventID, msg.Value, emailretry.BuildHeaders(meta)); pubErr != nil {
				span.End()
				return fmt.Errorf("publish to retry topic: %w", pubErr)
			}
			metrics.EmailRetried.Add(ctx, 1)
			span.End()
			sess.MarkMessage(msg, "")
			sess.Commit()
			continue
		}

		metrics.EmailSent.Add(ctx, 1)
		if h.inbox != nil {
			if err := h.inbox.MarkProcessed(ctx, h.consumerName, eventKey, reservationVersion); err != nil {
				status = "failed_inbox_complete"
				span.End()
				return err
			}
		}
		span.End()
		sess.MarkMessage(msg, "")
		sess.Commit()
	}
	return nil
}

// kafkaHeaderCarrier adapts sarama record headers to the OTel propagation
// HeaderCarrier interface so trace context can be extracted (Phase 5).
type kafkaHeaderCarrier []*sarama.RecordHeader

func (c kafkaHeaderCarrier) Get(key string) string {
	for _, h := range c {
		if h != nil && string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c kafkaHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for _, h := range c {
		if h != nil {
			keys = append(keys, string(h.Key))
		}
	}
	return keys
}

func (c kafkaHeaderCarrier) Set(string, string) {}

func (s *emailHandler) recordMetrics(method string, status string, start time.Time) {
	attrs := metric.WithAttributes(
		attribute.String("method", method),
		attribute.String("status", status),
	)
	s.requestCounter.Add(context.Background(), 1, attrs)
	s.requestDuration.Record(context.Background(), time.Since(start).Seconds(), attrs)
}
