package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type categoryStatsByIdError struct {
	logger logger.LoggerInterface
}

func NewCategoryStatsByIdError(logger logger.LoggerInterface) *categoryStatsByIdError {
	return &categoryStatsByIdError{logger: logger}
}

func (c *categoryStatsByIdError) HandleMonthTotalPriceError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CategoriesMonthlyTotalPriceResponse](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedFindMonthlyTotalPrice, fields...)
}

func (c *categoryStatsByIdError) HandleYearTotalPriceError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CategoriesYearlyTotalPriceResponse](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedFindYearlyTotalPrice, fields...)
}

func (c *categoryStatsByIdError) HandleMonthPrice(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CategoryMonthPriceResponse](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedFindMonthPrice, fields...)
}

func (c *categoryStatsByIdError) HandleYearPrice(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CategoryYearPriceResponse](c.logger, err, method, tracePrefix, span, status, category_errors.ErrFailedFindYearPrice, fields...)
}
