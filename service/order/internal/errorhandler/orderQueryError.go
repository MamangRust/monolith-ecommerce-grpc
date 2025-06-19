package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderQueryError struct {
	logger logger.LoggerInterface
}

func NewOrderQueryError(logger logger.LoggerInterface) *orderQueryError {
	return &orderQueryError{
		logger: logger,
	}
}

func (o *orderQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.OrderResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.OrderResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedFindAllOrders, fields...)
}

func (o *orderQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.OrderResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
