package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderStatsError struct {
	logger logger.LoggerInterface
}

func NewOrderStatsError(logger logger.LoggerInterface) *orderStatsError {
	return &orderStatsError{
		logger: logger,
	}
}

func (o *orderStatsError) HandleMonthlyTotalRevenueError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderMonthlyTotalRevenueResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindMonthlyTotalRevenue, fields...)
}

func (o *orderStatsError) HandleYearlyTotalRevenueError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderYearlyTotalRevenueResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindYearlyTotalRevenue, fields...)
}

func (o *orderStatsError) HandleMonthOrderStatsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderMonthlyResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderMonthlyResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindMonthlyOrder, fields...)
}

func (o *orderStatsError) HandleYearOrderStatsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderYearlyResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderYearlyResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindYearlyOrder, fields...)
}
