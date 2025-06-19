package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	CartQueryError   CartQueryError
	CartCommandError CartCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		CartQueryError:   NewCartQueryError(logger),
		CartCommandError: NewCartCommandError(logger),
	}
}
