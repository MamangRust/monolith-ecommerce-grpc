package server

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// PoolMetrics exposes PostgreSQL connection pool and Redis pool gauges using
// the global OTel meter provider. Values are refreshed periodically by the
// server monitoring task.
type PoolMetrics struct {
	pgTotalConns    metric.Int64Gauge
	pgIdleConns     metric.Int64Gauge
	pgAcquiredConns metric.Int64Gauge
	pgMaxConns      metric.Int64Gauge

	redisTotalConns metric.Int64Gauge
	redisIdleConns  metric.Int64Gauge
	redisHits       metric.Int64Gauge
	redisMisses     metric.Int64Gauge
}

func NewPoolMetrics(serviceName string) (*PoolMetrics, error) {
	meter := otel.Meter(serviceName)

	pgTotalConns, err := meter.Int64Gauge(
		"postgresql_pool_total_conns",
		metric.WithDescription("Total number of connections in the PostgreSQL pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	pgIdleConns, err := meter.Int64Gauge(
		"postgresql_pool_idle_conns",
		metric.WithDescription("Number of idle connections in the PostgreSQL pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	pgAcquiredConns, err := meter.Int64Gauge(
		"postgresql_pool_acquired_conns",
		metric.WithDescription("Number of currently acquired connections in the PostgreSQL pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	pgMaxConns, err := meter.Int64Gauge(
		"postgresql_pool_max_conns",
		metric.WithDescription("Maximum number of connections allowed in the PostgreSQL pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisTotalConns, err := meter.Int64Gauge(
		"redis_pool_total_conns",
		metric.WithDescription("Total number of connections in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisIdleConns, err := meter.Int64Gauge(
		"redis_pool_idle_conns",
		metric.WithDescription("Number of idle connections in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisHits, err := meter.Int64Gauge(
		"redis_pool_hits_total",
		metric.WithDescription("Total number of times a free connection was found in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	redisMisses, err := meter.Int64Gauge(
		"redis_pool_misses_total",
		metric.WithDescription("Total number of times a free connection was not found in the Redis pool"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &PoolMetrics{
		pgTotalConns:    pgTotalConns,
		pgIdleConns:     pgIdleConns,
		pgAcquiredConns: pgAcquiredConns,
		pgMaxConns:      pgMaxConns,
		redisTotalConns: redisTotalConns,
		redisIdleConns:  redisIdleConns,
		redisHits:       redisHits,
		redisMisses:     redisMisses,
	}, nil
}

func (m *PoolMetrics) RecordPGStats(ctx context.Context, total, idle, acquired, max int64) {
	m.pgTotalConns.Record(ctx, total)
	m.pgIdleConns.Record(ctx, idle)
	m.pgAcquiredConns.Record(ctx, acquired)
	m.pgMaxConns.Record(ctx, max)
}

func (m *PoolMetrics) RecordRedisStats(ctx context.Context, total, idle, hits, misses int64) {
	m.redisTotalConns.Record(ctx, total)
	m.redisIdleConns.Record(ctx, idle)
	m.redisHits.Record(ctx, hits)
	m.redisMisses.Record(ctx, misses)
}
