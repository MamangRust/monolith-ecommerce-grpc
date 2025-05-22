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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tokenService struct {
	trace        trace.Tracer
	refreshToken repository.RefreshTokenRepository
	token        auth.TokenManager
	logger       logger.LoggerInterface
}

func NewTokenService(refreshToken repository.RefreshTokenRepository, token auth.TokenManager, logger logger.LoggerInterface) *tokenService {
	return &tokenService{trace: otel.Tracer("token-service"), refreshToken: refreshToken, token: token, logger: logger}
}

func (s *tokenService) createAccessToken(ctx context.Context, id int) (string, error) {
	_, span := s.trace.Start(ctx, "TokenService.createAccessToken")
	defer span.End()

	span.SetAttributes(attribute.Int("user.id", id))

	s.logger.Debug("Creating access token",
		zap.Int("userID", id),
	)

	res, err := s.token.GenerateToken(id, "access")
	if err != nil {
		s.logger.Error("Failed to create access token",
			zap.Int("userID", id),
			zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Access token generation failed")
		return "", err
	}

	s.logger.Debug("Access token created successfully",
		zap.Int("userID", id),
	)
	span.SetStatus(codes.Ok, "Access token created")

	return res, nil
}

func (s *tokenService) createRefreshToken(ctx context.Context, id int) (string, error) {
	_, span := s.trace.Start(ctx, "TokenService.createRefreshToken")
	defer span.End()

	span.SetAttributes(attribute.Int("user.id", id))

	s.logger.Debug("Creating refresh token",
		zap.Int("userID", id),
	)

	res, err := s.token.GenerateToken(id, "refresh")
	if err != nil {
		traceID := traceunic.GenerateTraceID("REFRESH_TOKEN_FAILED")

		s.logger.Error("Failed to create refresh token",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Refresh token generation failed")
		return "", err
	}

	if err := s.refreshToken.DeleteRefreshTokenByUserId(id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		traceID := traceunic.GenerateTraceID("DELETE_REFRESH_TOKEN_ERR")
		s.logger.Error("Failed to delete existing refresh token", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete existing refresh token")
		return "", err
	}

	_, err = s.refreshToken.CreateRefreshToken(&requests.CreateRefreshToken{
		Token:     res,
		UserId:    id,
		ExpiresAt: time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("CREATE_REFRESH_TOKEN_ERR")

		s.logger.Error("Failed to create refresh token", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create refresh token")
		return "", err
	}

	s.logger.Debug("Refresh token created successfully",
		zap.Int("userID", id),
	)
	span.SetStatus(codes.Ok, "Refresh token created")

	return res, nil
}
