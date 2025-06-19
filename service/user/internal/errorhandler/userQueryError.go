package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userQueryError struct {
	logger logger.LoggerInterface
}

func NewUserQueryError(logger logger.LoggerInterface) *userQueryError {
	return &userQueryError{logger: logger}
}

func (u *userQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.UserResponse](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (u *userQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.UserResponseDeleteAt](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (u *userQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
