package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ProductQueryError interface {
	HandleRepositoryPaginationDetailError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
}

type ProductCommandError interface {
	HandleFileError(
		err error,
		method, tracePrefix, imagePath string,
		span trace.Span,
		status *string,
		fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleCreateProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponse, *response.ErrorResponse)
	HandleUpdateProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponse, *response.ErrorResponse)
	HandleTrashedProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
