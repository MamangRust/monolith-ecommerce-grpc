package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type TransactionStatsService interface {
	FindMonthlyAmountSuccess(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse)
	FindYearlyAmountSuccess(year int) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse)
	FindMonthlyAmountFailed(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse)
	FindYearlyAmountFailed(year int) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse)

	FindMonthlyMethodSuccess(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodSuccess(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)

	FindMonthlyMethodFailed(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodFailed(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)
}

type TransactionStatsByMerchantService interface {
	FindMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse)
	FindYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse)
	FindMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse)
	FindYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse)

	FindMonthlyMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)

	FindMonthlyMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)
}

type TransactionQueryService interface {
	FindAllTransactions(req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindByMerchant(req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse)
	FindByOrderId(orderID int) (*response.TransactionResponse, *response.ErrorResponse)
}

type TransactionCommandService interface {
	CreateTransaction(req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	UpdateTransaction(req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	TrashedTransaction(transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	RestoreTransaction(transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	DeleteTransactionPermanently(transactionID int) (bool, *response.ErrorResponse)
	RestoreAllTransactions() (bool, *response.ErrorResponse)
	DeleteAllTransactionPermanent() (bool, *response.ErrorResponse)
}
