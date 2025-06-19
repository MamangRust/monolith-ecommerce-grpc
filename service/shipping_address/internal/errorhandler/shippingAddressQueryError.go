package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type shippingAddressQueryError struct {
	logger logger.LoggerInterface
}

func NewShippingAddressQueryError(logger logger.LoggerInterface) *shippingAddressQueryError {
	return &shippingAddressQueryError{
		logger: logger,
	}
}

func (o *shippingAddressQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.ShippingAddressResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ShippingAddressResponse](o.logger, err, method, tracePrefix, span, status, shippingaddress_errors.ErrFailedFindAllShippingAddresses, fields...)
}

func (o *shippingAddressQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ShippingAddressResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (o *shippingAddressQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.ShippingAddressResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ShippingAddressResponse](o.logger, err, method, tracePrefix, span, status, shippingaddress_errors.ErrFailedFindShippingAddressByID, fields...)
}
