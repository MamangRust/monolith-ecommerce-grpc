package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MerchantAwardQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantAwardResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.MerchantAwardResponse, *response.ErrorResponse)
}

type MerchantAwardCommandError interface {
	HandleCreateMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantAwardResponse, *response.ErrorResponse)
	HandleUpdateMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantAwardResponse, *response.ErrorResponse)
	HandleTrashedMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.MerchantAwardResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllMerchantAwardError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
