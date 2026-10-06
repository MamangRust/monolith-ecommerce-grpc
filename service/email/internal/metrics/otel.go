package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	EmailSent       metric.Int64Counter
	EmailFailed     metric.Int64Counter
	EmailDuplicated metric.Int64Counter
	EmailInvalid    metric.Int64Counter
	EmailRetried    metric.Int64Counter
	EmailDeadLetter metric.Int64Counter
)

// Register creates all email service counters on the "email-service" meter.
// The instruments report through the OpenTelemetry SDK (OTLP) configured in main.
func Register() error {
	meter := otel.Meter("email-service")

	var err error
	if EmailSent, err = meter.Int64Counter("email_sent_total", metric.WithDescription("Total emails sent successfully")); err != nil {
		return err
	}
	if EmailFailed, err = meter.Int64Counter("email_failed_total", metric.WithDescription("Total emails failed")); err != nil {
		return err
	}
	if EmailDuplicated, err = meter.Int64Counter("email_duplicated_total", metric.WithDescription("Total duplicate Kafka messages skipped")); err != nil {
		return err
	}
	if EmailInvalid, err = meter.Int64Counter("email_invalid_total", metric.WithDescription("Total malformed or unprocessable Kafka messages")); err != nil {
		return err
	}
	if EmailRetried, err = meter.Int64Counter("email_retried_total", metric.WithDescription("Total emails published to the retry topic after a transient SMTP failure")); err != nil {
		return err
	}
	if EmailDeadLetter, err = meter.Int64Counter("email_deadletter_total", metric.WithDescription("Total emails dead-lettered after exhausting retries")); err != nil {
		return err
	}
	return nil
}
