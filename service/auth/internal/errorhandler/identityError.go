package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	refreshtoken_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/refresh_token_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type identityError struct {
	logger logger.LoggerInterface
}

func NewIdentityError(logger logger.LoggerInterface) *identityError {
	return &identityError{
		logger: logger,
	}
}

func (e *identityError) HandleInvalidTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedInValidToken,
		fields...,
	)
}

func (e *identityError) HandleExpiredRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedExpire,
		fields...,
	)
}

func HandleInvalidFormatUserIDError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorInvalidID[T](
		logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

func (e *identityError) HandleDeleteRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedDeleteRefreshToken,
		fields...,
	)
}

func (e *identityError) HandleUpdateRefreshTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedUpdateRefreshToken,
		fields...,
	)
}

func (e *identityError) HandleValidateTokenError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorTokenTemplate[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		refreshtoken_errors.ErrFailedInValidToken,
		fields...,
	)
}

func (e *identityError) HandleFindByIdError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

func (e *identityError) HandleGetMeError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}
