package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ShippingAddressQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.ShippingAddressResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.ShippingAddressResponse, *response.ErrorResponse)
}

type ShippingAddressCommandError interface {
	HandleTrashedShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllShippingAddressError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
