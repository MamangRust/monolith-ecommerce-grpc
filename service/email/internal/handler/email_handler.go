package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-ecommerce-email/internal/mailer"
	"github.com/MamangRust/monolith-ecommerce-email/internal/metrics"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// retryPublisher is implemented by *kafka.Kafka and lets the handlers publish
// retry/DLQ messages with metadata headers (Phase 4).
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
}

// kafkaHeaderCarrier adapts sarama record headers to the OTel propagation
// HeaderCarrier interface so trace context can be extracted.
type kafkaHeaderCarrier []*sarama.RecordHeader

func (c kafkaHeaderCarrier) Get(key string) string {
	for _, h := range c {
		if h != nil && string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c kafkaHeaderCarrier) Set(key, value string) {
	c = append(c, &sarama.RecordHeader{Key: []byte(key), Value: []byte(value)})
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

func NewEmailHandler(ctx context.Context, logger logger.LoggerInterface, mailer *mailer.Mailer) (*emailHandler, error) {
	meter := otel.Meter("email-service")

	requestCounter, err := meter.Int64Counter(
		"email_service_requests_total",
		metric.WithDescription("Total number of requests to the EmailService"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"email_service_request_duration_seconds",
		metric.WithDescription("Histogram of request durations for the EmailService"),
	)
	if err != nil {
		return nil, err
	}

	return &emailHandler{
		ctx:             ctx,
		logger:          logger,
		Mailer:          mailer,
		trace:           otel.Tracer("email-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}, nil
}

func (h *emailHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *emailHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *emailHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	start := time.Now()
	status := "success"

	defer func() {
		h.recordMetrics(h.ctx, "ConsumeClaim", status, start)
	}()

	_, span := h.trace.Start(h.ctx, "ConsumeClaim")
	defer span.End()

	for msg := range claim.Messages() {
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UNMARSHAL_MESSAGE")

			h.logger.Error("Failed to unmarshal message", zap.Error(err))

			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to unmarshal message")
			status = "failed_unmarshal_message"

			continue
		}

		email := payload["email"].(string)
		subject := payload["subject"].(string)
		body := payload["body"].(string)

		err := h.Mailer.Send(email, subject, body)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_SEND_EMAIL")

			h.logger.Error("Failed to send email", zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to send email")
			status = "failed_send_email"

			metrics.EmailFailed.Add(h.ctx, 1)
		} else {
			metrics.EmailSent.Add(h.ctx, 1)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func (s *emailHandler) recordMetrics(ctx context.Context, method string, status string, start time.Time) {
	attrs := metric.WithAttributes(attribute.String("method", method), attribute.String("status", status))
	s.requestCounter.Add(ctx, 1, attrs)
	s.requestDuration.Record(ctx, time.Since(start).Seconds(), attrs)
}
