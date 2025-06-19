package errorhandler

import "github.com/MamangRust/monolith-ecommerce-pkg/logger"

type ErrorHandler struct {
	TransactionQueryError          TransactionQueryError
	TransactionCommandError        TransactionCommandError
	TransactionStatsError          TransactionStatsError
	TransactonStatsByMerchantError TransactionStatsByMerchantError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		TransactionQueryError:          NewTransactionQueryError(logger),
		TransactionCommandError:        NewTransactionCommandError(logger),
		TransactionStatsError:          NewTransactionStatsError(logger),
		TransactonStatsByMerchantError: NewTransactionStatsByMerchantError(logger),
	}
}
