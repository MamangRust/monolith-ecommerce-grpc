package middleware

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// MetricsInterceptor records gRPC server request count, duration and errors
// per method and gRPC status code using the global OTel meter provider.
type MetricsInterceptor struct {
	requestCounter metric.Int64Counter
	durationHist   metric.Float64Histogram
	errorCounter   metric.Int64Counter
}

func NewMetricsInterceptor(serviceName string) (*MetricsInterceptor, error) {
	meter := otel.Meter(serviceName)

	requestCounter, err := meter.Int64Counter(
		"grpc_requests_total",
		metric.WithDescription("Total number of gRPC requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	durationHist, err := meter.Float64Histogram(
		"grpc_request_duration_seconds",
		metric.WithDescription("Duration of gRPC requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	errorCounter, err := meter.Int64Counter(
		"grpc_errors_total",
		metric.WithDescription("Total number of gRPC requests returning errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &MetricsInterceptor{
		requestCounter: requestCounter,
		durationHist:   durationHist,
		errorCounter:   errorCounter,
	}, nil
}

func (m *MetricsInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		code := status.Code(err)
		attrs := []attribute.KeyValue{
			attribute.String("method", info.FullMethod),
			attribute.String("grpc_status", code.String()),
		}

		m.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		m.durationHist.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(attrs...))

		if err != nil {
			m.errorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		}

		return resp, err
	}
}
