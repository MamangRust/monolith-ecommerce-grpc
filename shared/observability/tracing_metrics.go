package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TraceLoggerObservability interface {
	StartTracingAndLogging(
		ctx context.Context,
		method string,
		attrs ...attribute.KeyValue,
	) (
		context.Context,
		trace.Span,
		func(string),
		string,
		func(string, ...zap.Field),
	)
	RecordMetrics(ctx context.Context, method, status string, start time.Time)
}

type Observability struct {
	tracer          trace.Tracer
	meter           metric.Meter
	logger          logger.LoggerInterface
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	errorCounter    metric.Int64Counter
}

func NewObservability(serviceName string, logger logger.LoggerInterface) (TraceLoggerObservability, error) {
	tracer := otel.Tracer(serviceName)
	meter := otel.Meter(serviceName)

	requestCounter, err := meter.Int64Counter(
		"requests_total",
		metric.WithDescription("Total number of requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request counter: %w", err)
	}

	requestDuration, err := meter.Float64Histogram(
		"request_duration_seconds",
		metric.WithDescription("Duration of requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request duration histogram: %w", err)
	}

	errorCounter, err := meter.Int64Counter(
		"errors_total",
		metric.WithDescription("Total number of errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create error counter: %w", err)
	}

	return &Observability{
		tracer:          tracer,
		meter:           meter,
		logger:          logger,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
		errorCounter:    errorCounter,
	}, nil
}

func (o *Observability) StartTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := o.tracer.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	if o.logger != nil {
		o.logger.Info("Start: " + method)
	}

	end := func(status string) {
		o.RecordMetrics(ctx, method, status, start)

		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		if o.logger != nil {
			o.logger.Info(msg, fields...)
		}
	}

	return ctx, span, end, status, logSuccess
}

func (o *Observability) RecordMetrics(ctx context.Context, method, status string, start time.Time) {
	duration := time.Since(start).Seconds()

	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.String("status", status),
	}

	if o.requestCounter != nil {
		o.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}

	if o.requestDuration != nil {
		o.requestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
	}

	if status != "success" && o.errorCounter != nil {
		o.errorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}
