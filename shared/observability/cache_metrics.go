package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type CacheMetricsInterface interface {
	RecordCacheHit(ctx context.Context, key string)
	RecordCacheMiss(ctx context.Context, key string)
	RecordCacheSet(ctx context.Context, key string, success bool)
	RecordCacheDelete(ctx context.Context, key string, success bool)
	RecordCacheOperationLatency(ctx context.Context, operation string, duration time.Duration)
	RecordCacheError(ctx context.Context, operation, key string, err error)
}

type CacheMetrics struct {
	cacheHits         metric.Int64Counter
	cacheMisses       metric.Int64Counter
	cacheSets         metric.Int64Counter
	cacheSetErrors    metric.Int64Counter
	cacheDeletes      metric.Int64Counter
	cacheDeleteErrors metric.Int64Counter
	operationLatency  metric.Float64Histogram
	errors            metric.Int64Counter
}

func NewCacheMetrics(serviceName string) (CacheMetricsInterface, error) {
	meter := otel.Meter(serviceName)

	cacheHits, err := meter.Int64Counter(
		"cache_hits_total",
		metric.WithDescription("Total number of cache hits"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	cacheMisses, err := meter.Int64Counter(
		"cache_misses_total",
		metric.WithDescription("Total number of cache misses"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	cacheSets, err := meter.Int64Counter(
		"cache_sets_total",
		metric.WithDescription("Total number of cache set operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	cacheSetErrors, err := meter.Int64Counter(
		"cache_set_errors_total",
		metric.WithDescription("Total number of cache set errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	cacheDeletes, err := meter.Int64Counter(
		"cache_deletes_total",
		metric.WithDescription("Total number of cache delete operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	cacheDeleteErrors, err := meter.Int64Counter(
		"cache_delete_errors_total",
		metric.WithDescription("Total number of cache delete errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	operationLatency, err := meter.Float64Histogram(
		"cache_operation_duration_seconds",
		metric.WithDescription("Duration of cache operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	errors, err := meter.Int64Counter(
		"cache_errors_total",
		metric.WithDescription("Total number of cache errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &CacheMetrics{
		cacheHits:         cacheHits,
		cacheMisses:       cacheMisses,
		cacheSets:         cacheSets,
		cacheSetErrors:    cacheSetErrors,
		cacheDeletes:      cacheDeletes,
		cacheDeleteErrors: cacheDeleteErrors,
		operationLatency:  operationLatency,
		errors:            errors,
	}, nil
}

func (m *CacheMetrics) RecordCacheHit(ctx context.Context, key string) {
	attrs := []attribute.KeyValue{
		attribute.String("key", key),
	}
	m.cacheHits.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *CacheMetrics) RecordCacheMiss(ctx context.Context, key string) {
	attrs := []attribute.KeyValue{
		attribute.String("key", key),
	}
	m.cacheMisses.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *CacheMetrics) RecordCacheSet(ctx context.Context, key string, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("key", key),
		attribute.Bool("success", success),
	}
	m.cacheSets.Add(ctx, 1, metric.WithAttributes(attrs...))
	if !success {
		m.cacheSetErrors.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

func (m *CacheMetrics) RecordCacheDelete(ctx context.Context, key string, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("key", key),
		attribute.Bool("success", success),
	}
	m.cacheDeletes.Add(ctx, 1, metric.WithAttributes(attrs...))
	if !success {
		m.cacheDeleteErrors.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

func (m *CacheMetrics) RecordCacheOperationLatency(ctx context.Context, operation string, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
	}
	m.operationLatency.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (m *CacheMetrics) RecordCacheError(ctx context.Context, operation, key string, err error) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.String("key", key),
		attribute.String("error", err.Error()),
	}
	m.errors.Add(ctx, 1, metric.WithAttributes(attrs...))
}
