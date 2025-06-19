package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleCommandError struct {
	logger logger.LoggerInterface
}

func NewRoleCommandError(logger logger.LoggerInterface) *roleCommandError {
	return &roleCommandError{
		logger: logger,
	}
}

func (e *roleCommandError) HandleCreateRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedCreateRole,
		fields...,
	)
}

func (e *roleCommandError) HandleUpdateRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedUpdateRole,
		fields...,
	)
}

func (e *roleCommandError) HandleTrashedRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedTrashedRole,
		fields...,
	)
}

func (e *roleCommandError) HandleRestoreRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedRestoreRole,
		fields...,
	)
}

func (e *roleCommandError) HandleDeleteRolePermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedDeletePermanent,
		fields...,
	)
}

func (e *roleCommandError) HandleDeleteAllRolePermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedDeleteAll,
		fields...,
	)
}

func (e *roleCommandError) HandleRestoreAllRoleError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		role_errors.ErrFailedRestoreAll,
		fields...,
	)
}
