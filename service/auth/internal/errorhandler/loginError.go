package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type loginError struct {
	logger logger.LoggerInterface
}

func NewLoginError(logger logger.LoggerInterface) *loginError {
	return &loginError{
		logger: logger,
	}
}

func (e *loginError) HandleFindEmailError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TokenResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TokenResponse](
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
