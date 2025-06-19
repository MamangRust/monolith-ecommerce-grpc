package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	ReviewQueryError   ReviewQueryError
	ReviewCommandError ReviewCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		ReviewQueryError:   NewReviewQueryError(logger),
		ReviewCommandError: NewreviewCommandError(logger),
	}
}
