package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	refreshtoken_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/refresh_token_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type identityService struct {
	ctx             context.Context
	trace           trace.Tracer
	logger          logger.LoggerInterface
	token           auth.TokenManager
	refreshToken    repository.RefreshTokenRepository
	user            repository.UserRepository
	mapping         response_service.UserResponseMapper
	tokenService    tokenService
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewIdentityService(ctx context.Context, token auth.TokenManager, refreshToken repository.RefreshTokenRepository, user repository.UserRepository, logger logger.LoggerInterface, mapping response_service.UserResponseMapper, tokenService tokenService) *identityService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "identity_service_requests_total",
			Help: "Total number of auth requests",
		},
		[]string{"method", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "identity_service_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	return &identityService{
		ctx:             ctx,
		trace:           otel.Tracer("identity-service"),
		logger:          logger,
		token:           token,
		refreshToken:    refreshToken,
		user:            user,
		mapping:         mapping,
		tokenService:    tokenService,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *identityService) RefreshToken(token string) (*response.TokenResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("RefreshToken", status, startTime)
	}()

	ctx, span := s.trace.Start(s.ctx, "IdentityService.RefreshToken")
	defer span.End()

	span.SetAttributes(
		attribute.String("token", token),
	)

	s.logger.Debug("Refreshing token",
		zap.String("token", token),
	)

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) {
			traceID := traceunic.GenerateTraceID("TOKEN_EXPIRED")
			if err := s.refreshToken.DeleteRefreshToken(token); err != nil {
				s.logger.Error("Failed to delete expired refresh token",
					zap.String("trace_id", traceID),
					zap.Error(err),
				)
				span.RecordError(err)
				span.SetStatus(codes.Error, "Failed to delete expired token")
				status = "delete_token_failed"
				return nil, refreshtoken_errors.ErrFailedDeleteRefreshToken
			}

			s.logger.Error("Refresh token has expired",
				zap.String("trace_id", traceID),
				zap.Error(err),
			)
			span.RecordError(err)
			span.SetStatus(codes.Error, "Token expired")
			status = "token_expired"
			return nil, refreshtoken_errors.ErrFailedExpire
		}

		traceID := traceunic.GenerateTraceID("INVALID_TOKEN")
		s.logger.Error("Invalid refresh token",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid token")
		status = "invalid_token"
		return nil, refreshtoken_errors.ErrRefreshTokenNotFound
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		traceID := traceunic.GenerateTraceID("INVALID_USER_ID")
		s.logger.Error("Invalid user ID format in token",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid user ID format")
		status = "invalid_user_id"
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Invalid user ID format in token",
		}
	}

	span.SetAttributes(
		attribute.Int("user.id", userId),
	)

	accessToken, err := s.tokenService.createAccessToken(ctx, userId)
	if err != nil {
		traceID := traceunic.GenerateTraceID("ACCESS_TOKEN_FAILED")
		s.logger.Error("Failed to generate new access token",
			zap.String("trace_id", traceID),
			zap.Int("user_id", userId),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create access token")
		status = "access_token_failed"
		return nil, refreshtoken_errors.ErrFailedCreateAccess
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
	if err != nil {
		traceID := traceunic.GenerateTraceID("REFRESH_TOKEN_FAILED")
		s.logger.Error("Failed to generate new refresh token",
			zap.String("trace_id", traceID),
			zap.Int("user_id", userId),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create refresh token")
		status = "refresh_token_failed"
		return nil, refreshtoken_errors.ErrFailedCreateRefreshToken
	}

	expiryTime := time.Now().Add(24 * time.Hour)
	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshToken.UpdateRefreshToken(updateRequest); err != nil {
		traceID := traceunic.GenerateTraceID("UPDATE_TOKEN_FAILED")
		s.logger.Error("Failed to update refresh token in storage",
			zap.String("trace_id", traceID),
			zap.Int("user_id", userId),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update token")
		status = "token_update_failed"
		return nil, refreshtoken_errors.ErrFailedUpdateRefreshToken
	}

	s.logger.Debug("Refresh token refreshed successfully",
		zap.Int("user_id", userId),
	)
	span.SetStatus(codes.Ok, "Token refreshed successfully")

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *identityService) GetMe(token string) (*response.UserResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("GetMe", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "IdentityService.GetMe")
	defer span.End()

	span.SetAttributes(
		attribute.String("token", token),
	)

	s.logger.Debug("Fetching user details",
		zap.String("token", token),
	)

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		traceID := traceunic.GenerateTraceID("INVALID_ACCESS_TOKEN")
		s.logger.Error("Invalid access token",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid access token")
		status = "invalid_token"
		return nil, refreshtoken_errors.ErrFailedInValidToken
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		traceID := traceunic.GenerateTraceID("INVALID_USER_ID_FORMAT")
		s.logger.Error("Invalid user ID format in token",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid user ID format")
		status = "invalid_user_id"
		return nil, refreshtoken_errors.ErrFailedInValidUserId
	}

	span.SetAttributes(
		attribute.Int("user.id", userId),
	)

	user, err := s.user.FindById(userId)
	if err != nil {
		traceID := traceunic.GenerateTraceID("USER_NOT_FOUND")
		s.logger.Error("Failed to find user by ID",
			zap.String("trace_id", traceID),
			zap.Int("user_id", userId),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"
		return nil, user_errors.ErrUserNotFoundRes
	}

	userResponse := s.mapping.ToUserResponse(user)

	s.logger.Debug("User details fetched successfully",
		zap.Int("user_id", userId),
	)
	span.SetStatus(codes.Ok, "User details fetched")

	return userResponse, nil
}

func (s *identityService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
