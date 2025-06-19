package errorhandler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func handleErrorTemplate[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix, errorMessage string,
	span trace.Span,
	status *string,
	errorResp *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	traceID := traceunic.GenerateTraceID(tracePrefix)
	logMsg := fmt.Sprintf("%s in %s", errorMessage, method)

	allFields := append(fields,
		zap.Error(err),
		zap.String("trace.id", traceID),
	)

	logger.Error(logMsg, allFields...)

	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.AddEvent(logMsg)
	span.SetStatus(codes.Error, logMsg)

	*status = fmt.Sprintf("%s_error_%s", toSnakeCase(method), toSnakeCase(errorMessage))

	var zero T
	return zero, errorResp
}

func handleErrorRepository[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errorResp *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTemplate[T](
		logger, err, method, tracePrefix,
		"Repository error", span, status, errorResp, fields...,
	)
}

func handleErrorPagination[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errorResp *response.ErrorResponse,
	fields ...zap.Field,
) (T, *int, *response.ErrorResponse) {
	result, errResp := handleErrorRepository[T](
		logger, err, method, tracePrefix, span, status, errorResp, fields...,
	)
	return result, nil, errResp
}

func handleErrorTokenTemplate[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTemplate[T](logger, err, method, tracePrefix, "token", span, status, defaultErr, fields...)
}

func HandleTokenError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTokenTemplate[T](logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}

func handleErrorJSONMarshal[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTemplate[T](logger, err, method, tracePrefix, "json marshal", span, status, defaultErr, fields...)
}

func handleErrorKafkaSend[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTemplate[T](logger, err, method, tracePrefix, "kafka send", span, status, defaultErr, fields...)
}

func handleErrorGenerateRandomString[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTemplate[T](logger, err, method, tracePrefix, "generate random string", span, status, defaultErr, fields...)
}

func handleErrorInvalidID[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTemplate[T](logger, err, method, tracePrefix, "invalid id", span, status, defaultErr, fields...)
}

func handleErrorPasswordOperation[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix, operation string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorTemplate[T](logger, err, method, tracePrefix, operation, span, status, defaultErr, fields...)
}

func HandleRepositorySingleError[T any](
	logger logger.LoggerInterface,
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (T, *response.ErrorResponse) {
	return handleErrorRepository[T](logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}

func toSnakeCase(s string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
