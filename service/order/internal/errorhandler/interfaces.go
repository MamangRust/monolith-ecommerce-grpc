package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type OrderStatsError interface {
	HandleMonthlyTotalRevenueError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	HandleYearlyTotalRevenueError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)

	HandleMonthOrderStatsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderMonthlyResponse, *response.ErrorResponse)
	HandleYearOrderStatsError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderYearlyResponse, *response.ErrorResponse)
}

type OrderStatsByMerchantError interface {
	HandleMonthTotalRevenueByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse)
	HandleYearTotalRevenueByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse)

	HandleMonthOrderStatsByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderMonthlyResponse, *response.ErrorResponse)
	HandleYearOrderStatsByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.OrderYearlyResponse, *response.ErrorResponse)
}

type OrderQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.OrderResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
}

type OrderCommandError interface {
	HandleErrorInsufficientStockTemplate(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.OrderResponse, *response.ErrorResponse)
	HandleErrorInvalidCountStockTemplate(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.OrderResponse, *response.ErrorResponse)
	HandleCreateOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponse, *response.ErrorResponse)
	HandleUpdateOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponse, *response.ErrorResponse)
	HandleTrashedOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
