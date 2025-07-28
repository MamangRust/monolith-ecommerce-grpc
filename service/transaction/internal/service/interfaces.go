package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type TransactionStatsService interface {
	FindMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse)
	FindYearlyAmountSuccess(ctx context.Context, year int) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse)
	FindMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse)
	FindYearlyAmountFailed(ctx context.Context, year int) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse)

	FindMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodSuccess(ctx context.Context, year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)

	FindMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodFailed(ctx context.Context, year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)
}

type TransactionStatsByMerchantService interface {
	FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse)
	FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse)
	FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse)
	FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse)

	FindMonthlyMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)

	FindMonthlyMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse)
	FindYearlyMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse)
}

type TransactionQueryService interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, transactionID int) (*response.TransactionResponse, *response.ErrorResponse)
	FindByOrderId(ctx context.Context, orderID int) (*response.TransactionResponse, *response.ErrorResponse)
}

type TransactionCommandService interface {
	CreateTransaction(ctx context.Context, req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	UpdateTransaction(ctx context.Context, req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse)
	TrashedTransaction(ctx context.Context, transactionID int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	RestoreTransaction(ctx context.Context, transactionID int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse)
	DeleteTransactionPermanently(ctx context.Context, transactionID int) (bool, *response.ErrorResponse)
	RestoreAllTransactions(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
