package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewDetailCommandError struct {
	logger logger.LoggerInterface
}

func NewReviewDetailCommandError(logger logger.LoggerInterface) *reviewDetailCommandError {
	return &reviewDetailCommandError{
		logger: logger,
	}
}

func (o *reviewDetailCommandError) HandleInvalidFileError(
	err error,
	method, tracePrefix, imagePath string,
	span trace.Span,
	status *string,
	fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleFiledError(o.logger, err, method, tracePrefix, imagePath, span, status, fields...)
}

func (o *reviewDetailCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (o *reviewDetailCommandError) HandleCreateReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewDetailsResponse](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrFailedCreateReviewDetail, fields...)
}

func (o *reviewDetailCommandError) HandleUpdateReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewDetailsResponse](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrFailedUpdateReviewDetail, fields...)
}

func (o *reviewDetailCommandError) HandleTrashedReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewDetailsResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrFailedTrashedReviewDetail, fields...)
}

func (o *reviewDetailCommandError) HandleRestoreReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ReviewDetailsResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewDetailsResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrFailedRestoreReviewDetail, fields...)
}

func (o *reviewDetailCommandError) HandleDeleteReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrFailedDeletePermanentReview, fields...)
}

func (o *reviewDetailCommandError) HandleRestoreAllReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrFailedRestoreAllReviewDetail, fields...)
}

func (o *reviewDetailCommandError) HandleDeleteAllReviewDetailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrFailedDeleteAllReviewDetail, fields...)
}
