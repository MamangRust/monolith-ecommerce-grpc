package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	ReviewDetailQueryError   ReviewDetailQueryError
	ReviewDetailCommandError ReviewDetailCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		ReviewDetailQueryError:   NewReviewDetailQueryError(logger),
		ReviewDetailCommandError: NewReviewDetailCommandError(logger),
	}
}
