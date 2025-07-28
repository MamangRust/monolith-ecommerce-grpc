package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tokenService struct {
	refreshToken    repository.RefreshTokenRepository
	token           auth.TokenManager
	logger          logger.LoggerInterface
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewTokenService(
	refreshToken repository.RefreshTokenRepository, token auth.TokenManager, logger logger.LoggerInterface) *tokenService {
	prometheus.MustRegister(prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_service_requests_total",
			Help: "Total number of auth requests",
		},
		[]string{"method", "status"},
	))

	prometheus.MustRegister(prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "token_service_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	))

	return &tokenService{trace: otel.Tracer("token-service"), refreshToken: refreshToken, token: token, logger: logger}
}

func (s *tokenService) createAccessToken(ctx context.Context, id int) (string, error) {
	const method = "createAccessToken"

	ctx, end, logSuccess, status, logError := s.startTracingAndLogging(ctx, method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	res, err := s.token.GenerateToken(id, "access")
	if err != nil {
		status = "error"
		traceId := traceunic.GenerateTraceID("ACCESS_TOKEN_FAILED")

		logError(traceId, "Failed to create access token", err,
			zap.Int("userID", id),
			zap.Error(err),
		)

		return "", err
	}

	logSuccess("Created access token",
		zap.Int("userID", id),
	)

	return res, nil
}

func (s *tokenService) createRefreshToken(ctx context.Context, id int) (string, error) {
	const method = "createRefreshToken"

	ctx, end, logSuccess, status, logError := s.startTracingAndLogging(ctx, method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	res, err := s.token.GenerateToken(id, "refresh")
	if err != nil {
		status = "error"
		traceId := traceunic.GenerateTraceID("REFRESH_TOKEN_FAILED")

		logError(traceId, "Failed to create refresh token", err, zap.Int("user.id", id), zap.Error(err))
		return "", err
	}

	if err := s.refreshToken.DeleteRefreshTokenByUserId(ctx, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		status = "error"

		traceId := traceunic.GenerateTraceID("DELETE_REFRESH_TOKEN_ERR")

		logError(traceId, "Failed to delete existing refresh token", err, zap.Int("userID", id), zap.Error(err))

		return "", err
	}

	_, err = s.refreshToken.CreateRefreshToken(ctx, &requests.CreateRefreshToken{
		Token:     res,
		UserId:    id,
		ExpiresAt: time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		status = "error"

		traceId := traceunic.GenerateTraceID("CREATE_REFRESH_TOKEN_ERR")

		logError(traceId, "Failed to create refresh token", err, zap.Int("userID", id), zap.Error(err))

		return "", err
	}

	logSuccess("Created refresh token",
		zap.Int("userID", id),
	)

	return res, nil
}

func (s *tokenService) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (
	newCtx context.Context,
	end func(string),
	logSuccess func(string, ...zap.Field),
	status string,
	logError func(traceID string, msg string, err error, fields ...zap.Field),
) {
	start := time.Now()
	status = "success"

	newCtx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	end = func(status string) {
		s.recordMetrics(method, status, start)

		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}

		span.SetStatus(code, status)
		span.End()
	}

	logSuccess = func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError = func(traceID string, msg string, err error, fields ...zap.Field) {
		span.RecordError(err)
		span.SetStatus(codes.Error, msg)
		span.AddEvent(msg)

		allFields := append([]zap.Field{
			zap.String("trace.id", traceID),
			zap.Error(err),
		}, fields...)

		s.logger.Error(msg, allFields...)
	}

	return newCtx, end, logSuccess, status, logError
}

func (s *tokenService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
