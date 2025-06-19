package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CategoryStatsError interface {
	HandleMonthTotalPriceError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	HandleYearTotalPriceError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	HandleMonthPrice(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	HandleYearPrice(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryStatsByIdError interface {
	HandleMonthTotalPriceError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	HandleYearTotalPriceError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	HandleMonthPrice(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	HandleYearPrice(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryStatsByMerchantError interface {
	HandleMonthTotalPriceError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	HandleYearTotalPriceError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	HandleMonthPrice(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	HandleYearPrice(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CategoryResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.CategoryResponse, *response.ErrorResponse)
}

type CategoryCommandError interface {
	HandleCreateCategoryError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CategoryResponse, *response.ErrorResponse)
	HandleUpdateCategoryError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CategoryResponse, *response.ErrorResponse)
	HandleTrashedCategoryError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CategoryResponseDeleteAt, *response.ErrorResponse)
	HandleDeletePermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllPermanentlyError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CategoryResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreAllError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
