package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MerchantBusinessQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantBusinessResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.MerchantBusinessResponse, *response.ErrorResponse)
}

type MerchantBusinessCommandError interface {
	HandleCreateMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	HandleUpdateMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantBusinessResponse, *response.ErrorResponse)
	HandleTrashedMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantBusinessResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllMerchantBusinessError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
