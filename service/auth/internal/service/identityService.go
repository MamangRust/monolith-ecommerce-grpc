package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-auth/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type identityService struct {
	errorhandler    errorhandler.IdentityErrorHandler
	errorToken      errorhandler.TokenErrorHandler
	mencache        mencache.IdentityCache
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

func NewIdentityService(errohandler errorhandler.IdentityErrorHandler, errorToken errorhandler.TokenErrorHandler, mencache mencache.IdentityCache, token auth.TokenManager, refreshToken repository.RefreshTokenRepository, user repository.UserRepository, logger logger.LoggerInterface, mapping response_service.UserResponseMapper, tokenService tokenService) *identityService {
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
		errorhandler:    errohandler,
		errorToken:      errorToken,
		mencache:        mencache,
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

func (s *identityService) RefreshToken(ctx context.Context, token string) (*response.TokenResponse, *response.ErrorResponse) {
	const method = "RefreshToken"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("token", token))

	defer func() {
		end(status)
	}()

	if cachedUserID, found := s.mencache.GetRefreshToken(ctx, token); found {
		userId, err := strconv.Atoi(cachedUserID)
		if err == nil {
			s.mencache.DeleteRefreshToken(ctx, token)
			s.logger.Debug("Invalidated old refresh token from cache", zap.String("token", token))

			accessToken, err := s.tokenService.createAccessToken(ctx, userId)
			if err != nil {
				return s.errorToken.HandleCreateAccessTokenError(err, method, "CREATE_ACCESS_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
			}

			refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
			if err != nil {
				return s.errorToken.HandleCreateRefreshTokenError(err, method, "CREATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
			}

			expiryTime := time.Now().Add(24 * time.Hour)
			expirationDuration := time.Until(expiryTime)

			s.mencache.SetRefreshToken(ctx, refreshToken, expirationDuration)
			s.logger.Debug("Stored new refresh token in cache",
				zap.String("new_token", refreshToken),
				zap.Duration("expiration", expirationDuration))

			s.logger.Debug("Refresh token refreshed successfully (cached)", zap.Int("user_id", userId))
			span.SetStatus(codes.Ok, "Token refreshed successfully from cache")

			return &response.TokenResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}, nil
		}
	}

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) {
			s.mencache.DeleteRefreshToken(ctx, token)
			if err := s.refreshToken.DeleteRefreshToken(ctx, token); err != nil {

				return s.errorhandler.HandleDeleteRefreshTokenError(err, method, "DELETE_REFRESH_TOKEN", span, &status, zap.String("token", token))
			}

			return s.errorhandler.HandleExpiredRefreshTokenError(err, method, "TOKEN_EXPIRED", span, &status, zap.String("token", token))
		}

		return s.errorhandler.HandleInvalidTokenError(err, method, "INVALID_TOKEN", span, &status, zap.String("token", token))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {

		return errorhandler.HandleInvalidFormatUserIDError[*response.TokenResponse](s.logger, err, method, "INVALID_USER_ID", span, &status, zap.Int("user.id", userId))
	}

	span.SetAttributes(attribute.Int("user.id", userId))

	s.mencache.DeleteRefreshToken(ctx, token)
	if err := s.refreshToken.DeleteRefreshToken(ctx, token); err != nil {
		s.logger.Debug("Failed to delete old refresh token", zap.Error(err))

		return s.errorhandler.HandleDeleteRefreshTokenError(err, method, "DELETE_REFRESH_TOKEN", span, &status, zap.String("token", token))
	}

	accessToken, err := s.tokenService.createAccessToken(ctx, userId)
	if err != nil {

		return s.errorToken.HandleCreateAccessTokenError(err, method, "CREATE_ACCESS_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
	if err != nil {

		return s.errorToken.HandleCreateRefreshTokenError(err, method, "CREATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	expiryTime := time.Now().Add(24 * time.Hour)
	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshToken.UpdateRefreshToken(ctx, updateRequest); err != nil {
		s.mencache.DeleteRefreshToken(ctx, refreshToken)

		return s.errorhandler.HandleUpdateRefreshTokenError(err, method, "UPDATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	expirationDuration := time.Until(expiryTime)
	s.mencache.SetRefreshToken(ctx, refreshToken, expirationDuration)

	logSuccess("Refresh token refreshed successfully", zap.Int("user.id", userId))

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (s *identityService) GetMe(ctx context.Context, token string) (*response.UserResponse, *response.ErrorResponse) {
	const method = "GetMe"
	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("token", token))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Fetching user details", zap.String("token", token))

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleValidateTokenError(err, method, "INVALID_TOKEN", span, &status, zap.String("token", token))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		status = "error"

		return errorhandler.HandleInvalidFormatUserIDError[*response.UserResponse](
			s.logger, err, method, "INVALID_USER_ID", span, &status, zap.String("user_id_str", userIdStr),
		)
	}

	span.SetAttributes(attribute.Int("user.id", userId))

	if cachedUser, found := s.mencache.GetCachedUserInfo(ctx, userIdStr); found {
		logSuccess("User info retrieved from cache", zap.Int("user.id", userId))
		return cachedUser, nil
	}

	user, err := s.user.FindById(ctx, userId)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleFindByIdError(err, method, "FAILED_FETCH_USER", span, &status, zap.Int("user.id", userId))
	}

	userResponse := s.mapping.ToUserResponse(user)

	s.mencache.SetCachedUserInfo(ctx, userResponse, time.Minute*5)

	logSuccess("User details fetched successfully", zap.Int("user.id", userId))

	return userResponse, nil
}

func (s *identityService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *identityService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
