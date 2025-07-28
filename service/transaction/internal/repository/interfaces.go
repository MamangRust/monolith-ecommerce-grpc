package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*record.UserRecord, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, id int) (*record.MerchantRecord, error)
}

type OrderItemRepository interface {
	FindOrderItemByOrder(ctx context.Context, order_id int) ([]*record.OrderItemRecord, error)
}

type OrderQueryRepository interface {
	FindById(ctx context.Context, id int) (*record.OrderRecord, error)
}

type ShippingAddressQueryRepository interface {
	FindByOrder(ctx context.Context, order_id int) (*record.ShippingAddressRecord, error)
}

type TransactionStatsRepository interface {
	GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountSuccessRecord, error)
	GetYearlyAmountSuccess(ctx context.Context, year int) ([]*record.TransactionYearlyAmountSuccessRecord, error)
	GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountFailedRecord, error)
	GetYearlyAmountFailed(ctx context.Context, year int) ([]*record.TransactionYearlyAmountFailedRecord, error)

	GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*record.TransactionYearlyMethodRecord, error)

	GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*record.TransactionYearlyMethodRecord, error)
}

type TransactonStatsByMerchantRepository interface {
	GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountSuccessRecord, error)
	GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountSuccessRecord, error)
	GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountFailedRecord, error)
	GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountFailedRecord, error)

	GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error)

	GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error)
}

type TransactionQueryRepository interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*record.TransactionRecord, *int, error)
	FindById(ctx context.Context, transaction_id int) (*record.TransactionRecord, error)
	FindByOrderId(ctx context.Context, order_id int) (*record.TransactionRecord, error)
}

type TransactionCommandRepository interface {
	CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*record.TransactionRecord, error)
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error)
	TrashTransaction(ctx context.Context, transaction_id int) (*record.TransactionRecord, error)
	RestoreTransaction(ctx context.Context, transaction_id int) (*record.TransactionRecord, error)
	DeleteTransactionPermanently(ctx context.Context, transaction_id int) (bool, error)
	RestoreAllTransactions(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}
