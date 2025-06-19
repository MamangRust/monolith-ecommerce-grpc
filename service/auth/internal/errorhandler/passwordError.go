package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type passwordError struct {
	logger logger.LoggerInterface
}

func NewPasswordError(logger logger.LoggerInterface) *passwordError {
	return &passwordError{
		logger: logger,
	}
}

func (e *passwordError) HandlePasswordNotMatchError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorPasswordOperation[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		"not match",
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

func (e *passwordError) HandleHashPasswordError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorPasswordOperation[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		"hash",
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}

func (e *passwordError) HandleComparePasswordError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorPasswordOperation[*response.TokenResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		"compare",
		span,
		status,
		user_errors.ErrUserNotFoundRes,
		fields...,
	)
}
