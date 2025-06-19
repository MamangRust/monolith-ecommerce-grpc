package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDocumentQueryError struct {
	logger logger.LoggerInterface
}

func NewMerchantDocumentQueryError(logger logger.LoggerInterface) *merchantDocumentQueryError {
	return &merchantDocumentQueryError{
		logger: logger,
	}
}

func (e *merchantDocumentQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantDocumentResponse](e.logger, err, method, tracePrefix, span, status, merchantdocument_errors.ErrFailedFindAllMerchantDocuments, fields...)
}

func (e *merchantDocumentQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.MerchantDocumentResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
