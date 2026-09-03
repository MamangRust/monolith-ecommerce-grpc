package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// HTTPMetrics records HTTP request count, duration and errors per
// method, route pattern and status code using the global OTel meter provider.
type HTTPMetrics struct {
	requestCounter metric.Int64Counter
	durationHist   metric.Float64Histogram
	errorCounter   metric.Int64Counter
}

func NewHTTPMetrics(serviceName string) (*HTTPMetrics, error) {
	meter := otel.Meter(serviceName)

	requestCounter, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	durationHist, err := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	errorCounter, err := meter.Int64Counter(
		"http_errors_total",
		metric.WithDescription("Total number of HTTP requests returning 5xx errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &HTTPMetrics{
		requestCounter: requestCounter,
		durationHist:   durationHist,
		errorCounter:   errorCounter,
	}, nil
}

func (m *HTTPMetrics) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			status := c.Response().Status
			route := c.Path()
			if route == "" {
				route = "unmatched"
			}

			attrs := []attribute.KeyValue{
				attribute.String("method", c.Request().Method),
				attribute.String("route", route),
				attribute.Int("status", status),
			}

			m.requestCounter.Add(c.Request().Context(), 1, metric.WithAttributes(attrs...))
			m.durationHist.Record(c.Request().Context(), time.Since(start).Seconds(), metric.WithAttributes(attrs...))

			if status >= 500 {
				m.errorCounter.Add(c.Request().Context(), 1, metric.WithAttributes(attrs...))
			}

			return err
		}
	}
}
