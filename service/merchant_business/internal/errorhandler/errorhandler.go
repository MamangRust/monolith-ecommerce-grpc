package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	MerchantBusinessCommandError MerchantBusinessCommandError
	MerchantBusinessQueryError   MerchantBusinessQueryError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		MerchantBusinessCommandError: NewMerchantBusinessCommandError(logger),
		MerchantBusinessQueryError:   NewMerchantBusinessQueryError(logger),
	}
}
