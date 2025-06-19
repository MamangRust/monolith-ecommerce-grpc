package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type IdentityErrorHandler interface {
	HandleInvalidTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
	HandleExpiredRefreshTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
	HandleDeleteRefreshTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
	HandleUpdateRefreshTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
	HandleValidateTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleGetMeError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleFindByIdError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
}

type KafkaErrorHandler interface {
	HandleSendEmailForgotPassword(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleSendEmailRegister(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleSendEmailVerifyCode(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type LoginErrorHandler interface {
	HandleFindEmailError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
}

type MarshalErrorHandler interface {
	HandleMarshalRegisterError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleMarsalForgotPassword(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleMarshalVerifyCode(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type PasswordErrorHandler interface {
	HandlePasswordNotMatchError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleHashPasswordError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleComparePasswordError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
}

type PasswordResetErrorHandler interface {
	HandleFindEmailError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleCreateResetTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleFindTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleUpdatePasswordError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleDeleteTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleUpdateVerifiedError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleVerifyCodeError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type RandomStringErrorHandler interface {
	HandleRandomStringErrorRegister(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleRandomStringErrorForgotPassword(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

type RegisterErrorHandler interface {
	HandleAssignRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleFindRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleFindEmailError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
	HandleCreateUserError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
}

type TokenErrorHandler interface {
	HandleCreateAccessTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
	HandleCreateRefreshTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
}
