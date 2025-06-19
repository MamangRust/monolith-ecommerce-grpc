package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.UserResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
}

type UserCommandError interface {
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleCreateUserError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleUpdateUserError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleTrashedUserError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreUserError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteUserError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllUserError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllUserError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
