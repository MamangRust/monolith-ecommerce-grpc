package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/mailer"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/metrics"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type emailHandler struct {
	ctx             context.Context
	trace           trace.Tracer
	logger          logger.LoggerInterface
	Mailer          *mailer.Mailer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewEmailHandler(ctx context.Context, logger logger.LoggerInterface, mailer *mailer.Mailer) *emailHandler {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "email_service_requests_total",
			Help: "Total number of requests to the EmailService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "email_service_request_duration_seconds",
			Help:    "Histogram of request durations for the EmailService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &emailHandler{
		ctx:             ctx,
		logger:          logger,
		Mailer:          mailer,
		trace:           otel.Tracer("email-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
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

			metrics.EmailFailed.Inc()
		} else {
			metrics.EmailSent.Inc()
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func (s *emailHandler) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
