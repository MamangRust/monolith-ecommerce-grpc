package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	EmailSent       metric.Int64Counter
	EmailFailed     metric.Int64Counter
	EmailInvalid    metric.Int64Counter
	EmailRetried    metric.Int64Counter
	EmailDeadLetter metric.Int64Counter
)

// Register creates the OpenTelemetry counters from the global meter. Call it
// after the OTel SDK has been initialized in main (via pkg/otel, otel.go) so
// the counters are bound to the real meter provider and exported through the
// OTLP metric pipeline instead of a /metrics endpoint.
func Register() {
	meter := otel.Meter("email-service")

	EmailSent = mustCounter(meter, "email_sent_total", "Total emails sent successfully")
	EmailFailed = mustCounter(meter, "email_failed_total", "Total emails failed")
	EmailInvalid = mustCounter(meter, "email_invalid_total", "Total malformed or unprocessable Kafka messages")
	EmailRetried = mustCounter(meter, "email_retried_total", "Total emails published to the retry topic after a transient SMTP failure")
	EmailDeadLetter = mustCounter(meter, "email_deadletter_total", "Total emails dead-lettered after exhausting retries")
}

func mustCounter(meter metric.Meter, name, description string) metric.Int64Counter {
	counter, err := meter.Int64Counter(
		name,
		metric.WithDescription(description),
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(err)
	}
	return counter
}
