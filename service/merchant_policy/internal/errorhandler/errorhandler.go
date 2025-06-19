package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	MerchantPolicyCommandError MerchantPolicyCommandError
	MerchantPolicyQueryError   MerchantPolicyQueryError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		MerchantPolicyCommandError: NewMerchantPolicyCommandError(logger),
		MerchantPolicyQueryError:   NewMerchantPolicyQueryError(logger),
	}
}
