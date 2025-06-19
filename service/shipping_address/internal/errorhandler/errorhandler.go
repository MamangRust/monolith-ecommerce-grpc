package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	ShippingAddressQueryError   ShippingAddressQueryError
	ShippingAddressCommandError ShippingAddressCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		ShippingAddressQueryError:   NewShippingAddressQueryError(logger),
		ShippingAddressCommandError: NewShippingAddressCommandError(logger),
	}
}
