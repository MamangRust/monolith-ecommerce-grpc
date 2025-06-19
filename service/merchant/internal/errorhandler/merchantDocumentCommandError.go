package errorhandler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDocumentCommandError struct {
	logger logger.LoggerInterface
}

func NewMerchantDocumentCommandError(logger logger.LoggerInterface) *merchantDocumentCommandError {
	return &merchantDocumentCommandError{
		logger: logger,
	}
}

func (e *merchantDocumentCommandError) HandleCreateMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedCreateMerchantDocument,
		fields...,
	)
}

func (e *merchantDocumentCommandError) HandleUpdateMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedUpdateMerchantDocument,
		fields...,
	)
}

func (e *merchantDocumentCommandError) HandleUpdateMerchantDocumentStatusError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedUpdateMerchantDocument,
		fields...,
	)
}

func (e *merchantDocumentCommandError) HandleTrashedMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedTrashMerchantDocument,
		fields...,
	)
}

func (e *merchantDocumentCommandError) HandleRestoreMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.MerchantDocumentResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedRestoreMerchantDocument,
		fields...,
	)
}

func (e *merchantDocumentCommandError) HandleDeleteMerchantDocumentPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedDeleteMerchantDocument,
		fields...,
	)
}

func (e *merchantDocumentCommandError) HandleRestoreAllMerchantDocumentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedRestoreAllMerchantDocuments,
		fields...,
	)
}

func (e *merchantDocumentCommandError) HandleDeleteAllMerchantDocumentPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		merchantdocument_errors.ErrFailedDeleteAllMerchantDocuments,
		fields...,
	)
}
