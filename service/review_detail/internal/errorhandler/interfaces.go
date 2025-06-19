package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ReviewDetailQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.ReviewDetailsResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.ReviewDetailsResponse, *response.ErrorResponse)
}

type ReviewDetailCommandError interface {
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (bool, *response.ErrorResponse)
	HandleInvalidFileError(
		err error,
		method, tracePrefix, imagePath string,
		span trace.Span,
		status *string,
		fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleCreateReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	HandleUpdateReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewDetailsResponse, *response.ErrorResponse)
	HandleTrashedReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
