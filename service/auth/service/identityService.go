package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	mencache "github.com/MamangRust/monolith-ecommerce-auth/cache"
	"github.com/MamangRust/monolith-ecommerce-auth/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	sharederrorhandler "github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

// IdentityServiceDeps defines all dependencies required by IdentityService.
type IdentityServiceDeps struct {
	Cache         mencache.IdentityCache
	Token         auth.TokenManager
	RefreshToken  repository.RefreshTokenRepository
	User          repository.UserRepository
	Logger        logger.LoggerInterface
	TokenService  *tokenService
	Observability observability.TraceLoggerObservability
}

// identityService implements IdentityService.
type identityService struct {
	mencache      mencache.IdentityCache
	logger        logger.LoggerInterface
	token         auth.TokenManager
	refreshToken  repository.RefreshTokenRepository
	user          repository.UserRepository
	tokenService  *tokenService
	observability observability.TraceLoggerObservability
}

func NewIdentityService(param *IdentityServiceDeps) *identityService {
	return &identityService{
		mencache:      param.Cache,
		logger:        param.Logger,
		token:         param.Token,
		refreshToken:  param.RefreshToken,
		user:          param.User,
		tokenService:  param.TokenService,
		observability: param.Observability,
	}
}

func (s *identityService) RefreshToken(ctx context.Context, token string) (*response.TokenResponse, error) {
	const method = "RefreshToken"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("token", token))

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
				status = "error"
				return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", userId))
			}

			refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
			if err != nil {
				status = "error"
				return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", userId))
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
		status = "error"
		if errors.Is(err, auth.ErrTokenExpired) {
			s.mencache.DeleteRefreshToken(ctx, token)
			if err := s.refreshToken.DeleteRefreshToken(ctx, token); err != nil {
				return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.String("token", token))
			}
			expiredErr := fmt.Errorf("token expired: %w", err)
			return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, expiredErr, method, span, zap.String("token", token))
		}

		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.String("token", token))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.String("user_id_str", userIdStr))
	}

	span.SetAttributes(attribute.Int("user.id", userId))

	s.mencache.DeleteRefreshToken(ctx, token)
	if err := s.refreshToken.DeleteRefreshToken(ctx, token); err != nil {
		status = "error"
		s.logger.Debug("Failed to delete old refresh token", zap.Error(err))
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.String("token", token))
	}

	accessToken, err := s.tokenService.createAccessToken(ctx, userId)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", userId))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", userId))
	}

	expiryTime := time.Now().Add(24 * time.Hour)
	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshToken.UpdateRefreshToken(ctx, updateRequest); err != nil {
		status = "error"
		s.mencache.DeleteRefreshToken(ctx, refreshToken)
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", userId))
	}

	expirationDuration := time.Until(expiryTime)
	s.mencache.SetRefreshToken(ctx, refreshToken, expirationDuration)

	logSuccess("Refresh token refreshed successfully", zap.Int("user.id", userId))

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *identityService) GetMe(ctx context.Context, userId int) (*db.GetUserByIDRow, error) {
	const method = "GetMe"

	ctx, span, end, status, logSuccess :=
		s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	s.logger.Debug("Fetching user details", zap.Int("user.id", userId))

	cacheKey := strconv.Itoa(userId)

	if cachedUser, found := s.mencache.GetCachedUserInfo(ctx, cacheKey); found {
		logSuccess("User info retrieved from cache", zap.Int("user.id", userId))
		return cachedUser, nil
	}

	user, err := s.user.FindById(ctx, userId)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.GetUserByIDRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", userId),
		)
	}

	s.mencache.SetCachedUserInfo(ctx, user, 5*time.Minute)

	logSuccess("User details fetched successfully", zap.Int("user.id", userId))

	return user, nil
}
