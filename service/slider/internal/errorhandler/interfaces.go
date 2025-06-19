package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type SliderQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.SliderResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse)
}

type SliderCommandError interface {
	HandleCreateSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.SliderResponse, *response.ErrorResponse)
	HandleUpdateSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.SliderResponse, *response.ErrorResponse)
	HandleTrashedSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.SliderResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.SliderResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
