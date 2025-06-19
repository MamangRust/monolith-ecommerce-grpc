package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type MerchantQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
}

type MerchantDocumentQueryErrorHandler interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse)
}

type MerchantCommandErrorHandler interface {
	HandleCreateMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleUpdateMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleUpdateMerchantStatusError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleTrashedMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponseDeleteAt, *response.ErrorResponse)

	HandleRestoreMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantResponse, *response.ErrorResponse)

	HandleDeleteMerchantPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleRestoreAllMerchantError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleDeleteAllMerchantPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

type MerchantDocumentCommandErrorHandler interface {
	HandleCreateMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleUpdateMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleUpdateMerchantDocumentStatusError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleTrashedMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleRestoreMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	HandleDeleteMerchantDocumentPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleRestoreAllMerchantDocumentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	HandleDeleteAllMerchantDocumentPermanentError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}
