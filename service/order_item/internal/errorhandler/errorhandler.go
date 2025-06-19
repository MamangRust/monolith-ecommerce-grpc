package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	OrderItemQueryError OrderItemQueryError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		OrderItemQueryError: NewOrderItemQueryError(logger),
	}
}
