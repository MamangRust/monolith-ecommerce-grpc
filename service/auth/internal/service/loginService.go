package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	refreshtoken_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/refresh_token_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type loginService struct {
	ctx             context.Context
	logger          logger.LoggerInterface
	hash            hash.HashPassword
	user            repository.UserRepository
	refreshToken    repository.RefreshTokenRepository
	token           auth.TokenManager
	trace           trace.Tracer
	tokenService    tokenService
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewLoginService(
	ctx context.Context,
	logger logger.LoggerInterface,
	hash hash.HashPassword,
	userRepository repository.UserRepository,
	refreshToken repository.RefreshTokenRepository,
	token auth.TokenManager,
	tokenService tokenService,
) *loginService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_service_requests_total",
			Help: "Total number of auth requests",
		},
		[]string{"method", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "login_service_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &loginService{
		ctx:             ctx,
		logger:          logger,
		hash:            hash,
		user:            userRepository,
		refreshToken:    refreshToken,
		token:           token,
		trace:           otel.Tracer("login-service"),
		tokenService:    tokenService,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *loginService) Login(request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("Login", status, startTime)
	}()

	ctx, span := s.trace.Start(s.ctx, "LoginService.Login")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", request.Email),
	)

	s.logger.Debug("Starting login process",
		zap.String("email", request.Email),
	)

	res, err := s.user.FindByEmail(request.Email)
	if err != nil {
		traceID := traceunic.GenerateTraceID("EMAIL_NOT_FOUND")

		s.logger.Error("Failed to get user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"
		return nil, user_errors.ErrUserNotFoundRes
	}

	span.SetAttributes(
		attribute.Int("user.id", res.ID),
	)

	err = s.hash.ComparePassword(res.Password, request.Password)
	if err != nil {
		traceID := traceunic.GenerateTraceID("PASSWORD_MISMATCH")
		s.logger.Error("Failed to compare password", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Password mismatch")
		status = "password_mismatch"
		return nil, user_errors.ErrUserPassword
	}

	token, err := s.tokenService.createAccessToken(ctx, res.ID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("ACCESS_TOKEN_FAILED")

		s.logger.Error("Failed to generate JWT token", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create access token")
		status = "access_token_failed"
		return nil, refreshtoken_errors.ErrFailedCreateAccess
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, res.ID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("REFRESH_TOKEN_FAILED")
		s.logger.Error("Failed to generate refresh token", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create refresh token")
		status = "refresh_token_failed"
		return nil, refreshtoken_errors.ErrFailedCreateRefresh
	}

	s.logger.Debug("User logged in successfully",
		zap.String("email", request.Email),
		zap.Int("userID", res.ID),
	)
	span.SetStatus(codes.Ok, "Login successful")

	return &response.TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *loginService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
