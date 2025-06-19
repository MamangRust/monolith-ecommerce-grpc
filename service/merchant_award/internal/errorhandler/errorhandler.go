package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	MerchantAwardCommandError MerchantAwardCommandError
	MerchantAwardQueryError   MerchantAwardQueryError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		MerchantAwardCommandError: NewMerchantAwardCommandError(logger),
		MerchantAwardQueryError:   NewMerchantAwardQueryError(logger),
	}
}
