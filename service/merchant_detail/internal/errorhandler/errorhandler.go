package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	MerchantDetailCommandError     MerchantDetailCommandError
	MerchantDetailQueryError       MerchantDetailQueryError
	MerchantSocialLinkCommandError MerchantSocialLinkCommandError
	FileError                      FileError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		MerchantDetailCommandError:     NewMerchantDetailCommandError(logger),
		MerchantDetailQueryError:       NewMerchantDetailQueryError(logger),
		FileError:                      NewFileError(logger),
		MerchantSocialLinkCommandError: NewMerchantSocialLinkCommandError(logger),
	}
}
