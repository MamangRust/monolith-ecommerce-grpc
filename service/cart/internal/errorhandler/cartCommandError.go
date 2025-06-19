package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cartCommandError struct {
	logger logger.LoggerInterface
}

func NewCartCommandError(logger logger.LoggerInterface) *cartCommandError {
	return &cartCommandError{
		logger: logger,
	}
}
func (c *cartCommandError) HandleCreateCartError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CartResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CartResponse](c.logger, err, method, tracePrefix, span, status, cart_errors.ErrFailedCreateCart, fields...)
}

func (c *cartCommandError) HandleDeletePermanentError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, cart_errors.ErrFailedDeleteCart, fields...)
}

func (c *cartCommandError) HandleDeleteAllPermanentlyError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, cart_errors.ErrFailedDeleteAllCarts, fields...)
}
