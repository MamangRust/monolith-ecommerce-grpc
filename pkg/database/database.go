package database

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewClient(logger logger.LoggerInterface) (*pgxpool.Pool, error) {
	dbDriver := viper.GetString("DB_DRIVER")

	if dbDriver != "postgres" && dbDriver != "pgx" {
		logger.Error("pgxpool only supports PostgreSQL", zap.String("DB_DRIVER", dbDriver))
		return nil, fmt.Errorf("pgxpool only supports PostgreSQL, got: %s", dbDriver)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USERNAME"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PASSWORD"),
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Error("Failed to parse database config", zap.Error(err))
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	maxOpenConns := viper.GetInt("DB_MAX_OPEN_CONNS")
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}
	config.MaxConns = int32(maxOpenConns)

	minIdleConns := viper.GetInt("DB_MIN_IDLE_CONNS")
	if minIdleConns <= 0 {
		minIdleConns = 50
	}
	config.MinConns = int32(minIdleConns)

	connMaxLifetime := viper.GetDuration("DB_CONN_MAX_LIFETIME")
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour
	}
	config.MaxConnLifetime = connMaxLifetime

	connMaxIdleTime := viper.GetDuration("DB_CONN_MAX_IDLE_TIME")
	if connMaxIdleTime == 0 {
		connMaxIdleTime = 30 * time.Minute
	}
	config.MaxConnIdleTime = connMaxIdleTime

	healthCheckPeriod := viper.GetDuration("DB_HEALTH_CHECK_PERIOD")
	if healthCheckPeriod == 0 {
		healthCheckPeriod = time.Minute
	}
	config.HealthCheckPeriod = healthCheckPeriod

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("Failed to create connection pool", zap.Error(err))
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Debug("Database connection pool established successfully",
		zap.String("DB_DRIVER", "pgx"),
		zap.Int32("MaxConns", config.MaxConns),
		zap.Int32("MinConns", config.MinConns),
		zap.Duration("MaxConnLifetime", config.MaxConnLifetime),
		zap.Duration("MaxConnIdleTime", config.MaxConnIdleTime),
		zap.Duration("HealthCheckPeriod", config.HealthCheckPeriod),
	)

	return pool, nil
}
