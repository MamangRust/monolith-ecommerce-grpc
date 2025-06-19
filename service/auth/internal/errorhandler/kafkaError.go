package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type kafkaError struct {
	logger logger.LoggerInterface
}

func NewKafkaError(logger logger.LoggerInterface) *kafkaError {
	return &kafkaError{
		logger: logger,
	}
}

func (e *kafkaError) HandleSendEmailForgotPassword(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorKafkaSend[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrFailedSendEmail,
		fields...,
	)
}

func (e *kafkaError) HandleSendEmailRegister(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorKafkaSend[*response.UserResponse](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrFailedSendEmail,
		fields...,
	)
}

func (e *kafkaError) HandleSendEmailVerifyCode(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorKafkaSend[bool](
		e.logger,
		err,
		method,
		tracePrefix,
		span,
		status,
		user_errors.ErrFailedSendEmail,
		fields...,
	)
}
