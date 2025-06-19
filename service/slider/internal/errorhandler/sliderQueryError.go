package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type sliderQueryError struct {
	logger logger.LoggerInterface
}

func NewSliderQueryError(logger logger.LoggerInterface) *sliderQueryError {
	return &sliderQueryError{
		logger: logger,
	}
}

func (o *sliderQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.SliderResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.SliderResponse](o.logger, err, method, tracePrefix, span, status, slider_errors.ErrFailedFindAllSliders, fields...)
}

func (o *sliderQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.SliderResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
