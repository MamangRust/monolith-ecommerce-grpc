package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	MerchantDetailCommandError MerchantDetailCommandError
	MerchantDetailQueryError   MerchantDetailQueryError
	FileError                  FileError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		MerchantDetailCommandError: NewMerchantDetailCommandError(logger),
		MerchantDetailQueryError:   NewMerchantDetailQueryError(logger),
		FileError:                  NewFileError(logger),
	}
}
