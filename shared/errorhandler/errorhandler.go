package errorhandler

import (
	"net/http"

	errorsstd "errors"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func HandleError[T any](
	logger logger.LoggerInterface,
	err error,
	method string,
	span trace.Span,
	fields ...zap.Field,
) (T, error) {
	var zero T
	traceID := span.SpanContext().TraceID().String()

	var appErr *errors.AppError
	statusCode := http.StatusInternalServerError
	var internalErr error

	if errorsstd.As(err, &appErr) {
		statusCode = appErr.Code
		internalErr = appErr.Internal
	}

	logFields := append([]zap.Field{
		zap.String("method", method),
		zap.String("trace.id", traceID),
		zap.Int("status_code", statusCode),
	}, fields...)

	if internalErr != nil {
		logFields = append(logFields, zap.Error(internalErr))
	} else {
		logFields = append(logFields, zap.Error(err))
	}

	logger.Error("request failed", logFields...)

	span.SetAttributes(
		attribute.String("trace.id", traceID),
		attribute.Int("http.status_code", statusCode),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return zero, err
}


func HandlerErrorPagination[T any](
	logger logger.LoggerInterface,
	err error,
	method string,
	span trace.Span,
	fields ...zap.Field,
) (T, *int, error) {
	result, err := HandleError[T](logger, err, method, span, fields...)
	return result, nil, err
}
