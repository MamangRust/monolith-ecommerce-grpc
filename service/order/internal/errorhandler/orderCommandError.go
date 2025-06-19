package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderCommandError struct {
	logger logger.LoggerInterface
}

func NewOrderCommandError(logger logger.LoggerInterface) *orderCommandError {
	return &orderCommandError{
		logger: logger,
	}
}

func (o *orderCommandError) HandleErrorInsufficientStockTemplate(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.OrderResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.OrderResponse](o.logger, err, method, tracePrefix, "Insufficient stock", span, status, errResp, fields...)
}

func (o *orderCommandError) HandleErrorInvalidCountStockTemplate(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.OrderResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.OrderResponse](o.logger, err, method, tracePrefix, "Invalid count stock", span, status, errResp, fields...)
}

func (o *orderCommandError) HandleCreateOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.OrderResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedCreateOrder, fields...)
}

func (o *orderCommandError) HandleUpdateOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.OrderResponse](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedUpdateOrder, fields...)
}

func (o *orderCommandError) HandleTrashedOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.OrderResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedTrashOrder, fields...)
}

func (o *orderCommandError) HandleRestoreOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.OrderResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedRestoreOrder, fields...)
}

func (o *orderCommandError) HandleDeleteOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedDeleteOrderPermanent, fields...)
}

func (o *orderCommandError) HandleRestoreAllOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedRestoreAllOrder, fields...)
}

func (o *orderCommandError) HandleDeleteAllOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, order_errors.ErrFailedDeleteAllOrderPermanent, fields...)
}
