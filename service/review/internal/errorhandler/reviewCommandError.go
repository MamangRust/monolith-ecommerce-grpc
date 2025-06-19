package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewCommandError struct {
	logger logger.LoggerInterface
}

func NewreviewCommandError(logger logger.LoggerInterface) *reviewCommandError {
	return &reviewCommandError{
		logger: logger,
	}
}

func (o *reviewCommandError) HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.ReviewResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewResponse](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (o *reviewCommandError) HandleCreateReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewResponse](o.logger, err, method, tracePrefix, span, status, review_errors.ErrFailedCreateReview, fields...)
}

func (o *reviewCommandError) HandleUpdateReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewResponse](o.logger, err, method, tracePrefix, span, status, review_errors.ErrFailedUpdateReview, fields...)
}

func (o *reviewCommandError) HandleTrashedReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, review_errors.ErrFailedTrashedReview, fields...)
}

func (o *reviewCommandError) HandleRestoreReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, review_errors.ErrFailedRestoreReview, fields...)
}

func (o *reviewCommandError) HandleDeleteReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, review_errors.ErrFailedDeletePermanentReview, fields...)
}

func (o *reviewCommandError) HandleRestoreAllReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, review_errors.ErrFailedRestoreAllReviews, fields...)
}

func (o *reviewCommandError) HandleDeleteAllReviewError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, review_errors.ErrFailedDeleteAllPermanentReviews, fields...)
}
