package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewQueryError struct {
	logger logger.LoggerInterface
}

func NewReviewQueryError(logger logger.LoggerInterface) *reviewQueryError {
	return &reviewQueryError{
		logger: logger,
	}
}

func (o *reviewQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.ReviewResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ReviewResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindAllOrders, fields...)
}

func (o *reviewQueryError) HandleRepositoryPaginationDetailError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.ReviewsDetailResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ReviewsDetailResponse](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (o *reviewQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.ReviewResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ReviewResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (o *reviewQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.ReviewResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ReviewResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindOrderById, fields...)
}
