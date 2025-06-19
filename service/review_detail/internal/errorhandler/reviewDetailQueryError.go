package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewDetailQueryError struct {
	logger logger.LoggerInterface
}

func NewReviewDetailQueryError(logger logger.LoggerInterface) *reviewDetailQueryError {
	return &reviewDetailQueryError{
		logger: logger,
	}
}

func (o *reviewDetailQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.ReviewDetailsResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ReviewDetailsResponse](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrFailedFindAllReview, fields...)
}

func (o *reviewDetailQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.ReviewDetailsResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ReviewDetailsResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (o *reviewDetailQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.ReviewDetailsResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewDetailsResponse](o.logger, err, method, tracePrefix, span, status, reviewdetail_errors.ErrReviewDetailNotFoundRes, fields...)
}
