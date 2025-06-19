package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type RoleCommandErrorHandler interface {
	HandleCreateRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
	HandleUpdateRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
	HandleTrashedRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
	HandleRestoreRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
	HandleDeleteRolePermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleDeleteAllRolePermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleRestoreAllRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type RoleQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.RoleResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeletedError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.RoleResponse, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		defaultErr *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.RoleResponse, *response.ErrorResponse)
}
