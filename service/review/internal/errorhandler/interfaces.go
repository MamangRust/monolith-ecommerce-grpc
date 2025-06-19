package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ReviewQueryError interface {
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
	) ([]*response.ReviewResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.ReviewResponse, *response.ErrorResponse)
}

type ReviewCommandError interface {
	HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.ReviewResponse, *response.ErrorResponse)
	HandleCreateReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewResponse, *response.ErrorResponse)
	HandleUpdateReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewResponse, *response.ErrorResponse)
	HandleTrashedReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
