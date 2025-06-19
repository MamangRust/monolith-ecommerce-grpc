package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type shippingAddressCommandError struct {
	logger logger.LoggerInterface
}

func NewShippingAddressCommandError(logger logger.LoggerInterface) *shippingAddressCommandError {
	return &shippingAddressCommandError{
		logger: logger,
	}
}
func (o *shippingAddressCommandError) HandleTrashedShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ShippingAddressResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, shippingaddress_errors.ErrFailedTrashShippingAddress, fields...)
}

func (o *shippingAddressCommandError) HandleRestoreShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ShippingAddressResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, shippingaddress_errors.ErrFailedRestoreShippingAddress, fields...)
}

func (o *shippingAddressCommandError) HandleDeleteShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, shippingaddress_errors.ErrFailedDeleteShippingAddressPermanent, fields...)
}

func (o *shippingAddressCommandError) HandleRestoreAllShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, shippingaddress_errors.ErrFailedRestoreAllShippingAddresses, fields...)
}

func (o *shippingAddressCommandError) HandleDeleteAllShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](o.logger, err, method, tracePrefix, span, status, shippingaddress_errors.ErrFailedDeleteAllShippingAddressesPermanent, fields...)
}
