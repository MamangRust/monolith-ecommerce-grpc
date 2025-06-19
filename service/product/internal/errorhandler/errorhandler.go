package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	ProductQueryError   ProductQueryError
	ProductCommandError ProductCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		ProductQueryError:   NewProductQueryError(logger),
		ProductCommandError: NewProductCommandError(logger),
	}
}
