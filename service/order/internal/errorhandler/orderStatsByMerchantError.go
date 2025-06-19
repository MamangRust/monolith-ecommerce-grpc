package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderStatsByMerchantError struct {
	logger logger.LoggerInterface
}

func NewOrderStatsByMerchantError(logger logger.LoggerInterface) *orderStatsByMerchantError {
	return &orderStatsByMerchantError{
		logger: logger,
	}
}

func (o *orderStatsByMerchantError) HandleMonthTotalRevenueByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderMonthlyTotalRevenueResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindMonthlyTotalRevenueByMerchant, fields...)
}

func (o *orderStatsByMerchantError) HandleYearTotalRevenueByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderYearlyTotalRevenueResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindYearlyTotalRevenueByMerchant, fields...)
}

func (o *orderStatsByMerchantError) HandleMonthOrderStatsByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderMonthlyResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderMonthlyResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindMonthlyOrderByMerchant, fields...)
}

func (o *orderStatsByMerchantError) HandleYearOrderStatsByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderYearlyResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderYearlyResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindYearlyOrderByMerchant, fields...)
}
