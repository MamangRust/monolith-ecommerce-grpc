package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	refreshtoken_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/refresh_token_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tokenError struct {
	logger logger.LoggerInterface
}

func NewTokenError(logger logger.LoggerInterface) *tokenError {
	return &tokenError{
		logger: logger,
	}
}

func (e *tokenError) HandleCreateAccessTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedCreateAccess,
		fields...,
	)
}

func (e *tokenError) HandleCreateRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedCreateRefresh,
		fields...,
	)
}
