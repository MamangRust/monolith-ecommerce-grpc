package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type tokenService struct {
	refreshToken  repository.RefreshTokenRepository
	token         auth.TokenManager
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

// tokenServiceDeps holds the dependencies required to construct a tokenService.
type tokenServiceDeps struct {
	Token         auth.TokenManager
	RefreshToken  repository.RefreshTokenRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewTokenService(
	params *tokenServiceDeps,
) *tokenService {
	return &tokenService{
		refreshToken:  params.RefreshToken,
		token:         params.Token,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *tokenService) createAccessToken(ctx context.Context, id int) (string, error) {
	const method = "createAccessToken"

	ctx, _, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	res, err := s.token.GenerateToken(id, "access")
	if err != nil {
		status = "error"
		traceId := traceunic.GenerateTraceID("ACCESS_TOKEN_FAILED")

		s.logger.Error("Failed to create access token",
			zap.String("traceId", traceId),
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

	ctx, _, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	res, err := s.token.GenerateToken(id, "refresh")
	if err != nil {
		status = "error"
		traceId := traceunic.GenerateTraceID("REFRESH_TOKEN_FAILED")

		s.logger.Error("Failed to create refresh token",
			zap.String("traceId", traceId),
			zap.Int("userID", id),
			zap.Error(err),
		)
		return "", err
	}

	if err := s.refreshToken.DeleteRefreshTokenByUserId(ctx, id); err != nil {
		status = "error"
		traceId := traceunic.GenerateTraceID("DELETE_REFRESH_TOKEN_ERR")

		s.logger.Error("Failed to delete existing refresh token",
			zap.String("traceId", traceId),
			zap.Error(err),
			zap.Int("userID", id),
		)

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

		s.logger.Error("Failed to create refresh token",
			zap.String("traceId", traceId),
			zap.Error(err),
			zap.Int("userID", id),
		)

		return "", err
	}

	logSuccess("Created refresh token",
		zap.Int("userID", id),
	)

	return res, nil
}
