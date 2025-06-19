package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleQueryError struct {
	logger logger.LoggerInterface
}

func NewRoleQueryError(logger logger.LoggerInterface) *roleQueryError {
	return &roleQueryError{
		logger: logger,
	}
}

func (e *roleQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.RoleResponse](
		e.logger, err, method, tracePrefix, span, status, role_errors.ErrFailedFindAll, fields...,
	)
}

func (e *roleQueryError) HandleRepositoryPaginationDeletedError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.RoleResponseDeleteAt](
		e.logger, err, method, tracePrefix, span, status, errResp, fields...,
	)
}

func (e *roleQueryError) HandleRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.RoleResponse](e.logger, err, method, tracePrefix, span, status, role_errors.ErrFailedFindAll, fields...)
}

func (e *roleQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](e.logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}
