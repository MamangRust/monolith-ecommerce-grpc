package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	SliderQueryError   SliderQueryError
	SliderCommandError SliderCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		SliderQueryError:   NewSliderQueryError(logger),
		SliderCommandError: NewSliderCommandError(logger),
	}
}
