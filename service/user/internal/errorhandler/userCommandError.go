package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userCommandError struct {
	logger logger.LoggerInterface
}

func NewUserCommandError(logger logger.LoggerInterface) *userCommandError {
	return &userCommandError{logger: logger}
}
func (u *userCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (u *userCommandError) HandleCreateUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (u *userCommandError) HandleUpdateUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (u *userCommandError) HandleTrashedUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponseDeleteAt](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (u *userCommandError) HandleRestoreUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponseDeleteAt](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (u *userCommandError) HandleDeleteUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (u *userCommandError) HandleRestoreAllUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (u *userCommandError) HandleDeleteAllUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}
