package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	OrderCommandError    OrderCommandError
	OrderQueryError      OrderQueryError
	OrderStats           OrderStatsError
	OrderStatsByMerchant OrderStatsByMerchantError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		OrderCommandError:    NewOrderCommandError(logger),
		OrderQueryError:      NewOrderQueryError(logger),
		OrderStats:           NewOrderStatsError(logger),
		OrderStatsByMerchant: NewOrderStatsByMerchantError(logger),
	}
}
