package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type sliderCommandError struct {
	logger logger.LoggerInterface
}

func NewSliderCommandError(logger logger.LoggerInterface) *sliderCommandError {
	return &sliderCommandError{
		logger: logger,
	}
}

func (o *sliderCommandError) HandleCreateSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.SliderResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.SliderResponse](o.logger, err, method, tracePrefix, span, status, slider_errors.ErrFailedCreateSlider, fields...)
}

func (o *sliderCommandError) HandleUpdateSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.SliderResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.SliderResponse](o.logger, err, method, tracePrefix, span, status, slider_errors.ErrFailedUpdateSlider, fields...)
}

func (o *sliderCommandError) HandleTrashedSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.SliderResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, slider_errors.ErrFailedTrashSlider, fields...)
}

func (o *sliderCommandError) HandleRestoreSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.SliderResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, slider_errors.ErrFailedRestoreSlider, fields...)
}

func (o *sliderCommandError) HandleDeleteSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, slider_errors.ErrFailedDeletePermanentSlider, fields...)
}

func (o *sliderCommandError) HandleRestoreAllSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, slider_errors.ErrFailedRestoreAllSliders, fields...)
}

func (o *sliderCommandError) HandleDeleteAllSliderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, slider_errors.ErrFailedDeleteAllPermanentSliders, fields...)
}
