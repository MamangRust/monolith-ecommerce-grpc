package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type TransactionStatsService interface {
	FindMonthlyAmountSuccess(
		ctx context.Context,
		req *requests.MonthAmountTransaction,
	) ([]*db.GetMonthlyAmountTransactionSuccessRow, error)

	FindYearlyAmountSuccess(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyAmountTransactionSuccessRow, error)

	FindMonthlyAmountFailed(
		ctx context.Context,
		req *requests.MonthAmountTransaction,
	) ([]*db.GetMonthlyAmountTransactionFailedRow, error)

	FindYearlyAmountFailed(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyAmountTransactionFailedRow, error)

	FindMonthlyMethodSuccess(
		ctx context.Context,
		req *requests.MonthMethodTransaction,
	) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error)

	FindYearlyMethodSuccess(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTransactionMethodsSuccessRow, error)

	FindMonthlyMethodFailed(
		ctx context.Context,
		req *requests.MonthMethodTransaction,
	) ([]*db.GetMonthlyTransactionMethodsFailedRow, error)

	FindYearlyMethodFailed(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTransactionMethodsFailedRow, error)
}

type TransactionStatsByMerchantService interface {
	FindMonthlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error)

	FindYearlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error)

	FindMonthlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error)

	FindYearlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error)

	FindMonthlyMethodByMerchantSuccess(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error)

	FindYearlyMethodByMerchantSuccess(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error)

	FindMonthlyMethodByMerchantFailed(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error)

	FindYearlyMethodByMerchantFailed(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error)
}

type TransactionQueryService interface {
	FindAllTransactions(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsTrashedRow, *int, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllTransactionByMerchant,
	) ([]*db.GetTransactionByMerchantRow, *int, error)

	FindById(
		ctx context.Context,
		transaction_id int,
	) (*db.GetTransactionByIDRow, error)

	FindByOrderId(
		ctx context.Context,
		order_id int,
	) (*db.GetTransactionByOrderIDRow, error)
}

type TransactionCommandService interface {
	CreateTransaction(
		ctx context.Context,
		request *requests.CreateTransactionRequest,
	) (*db.CreateTransactionRow, error)

	UpdateTransaction(
		ctx context.Context,
		request *requests.UpdateTransactionRequest,
	) (*db.UpdateTransactionRow, error)

	TrashedTransaction(
		ctx context.Context,
		transaction_id int,
	) (*db.Transaction, error)

	RestoreTransaction(
		ctx context.Context,
		transaction_id int,
	) (*db.Transaction, error)

	DeleteTransactionPermanently(
		ctx context.Context,
		transaction_id int,
	) (bool, error)

	RestoreAllTransactions(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}
