package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	MerchantQueryError           MerchantQueryErrorHandler
	MerchantCommandError         MerchantCommandErrorHandler
	MerchantDocumentQueryError   MerchantDocumentQueryErrorHandler
	MerchantDocumentCommandError MerchantDocumentCommandErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		MerchantQueryError: NewMerchantQueryError(logger),
		MerchantCommandError: NewMerchantCommandError(
			logger,
		),
		MerchantDocumentQueryError: NewMerchantDocumentQueryError(
			logger,
		),
		MerchantDocumentCommandError: NewMerchantDocumentCommandError(
			logger),
	}
}
