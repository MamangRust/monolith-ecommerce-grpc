package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
}

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type OrderItemRepository interface {
	FindOrderItemByOrder(
		ctx context.Context,
		order_id int,
	) ([]*db.GetOrderItemsByOrderRow, error)
}

type OrderQueryRepository interface {
	FindByID(
		ctx context.Context,
		order_id int,
	) (*db.GetOrderByIDRow, error)
}

type ShippingAddressQueryRepository interface {
	FindByID(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingAddressByOrderIDRow, error)
}

type TransactionStatsRepository interface {
	GetMonthlyAmountSuccess(
		ctx context.Context,
		req *requests.MonthAmountTransaction,
	) ([]*db.GetMonthlyAmountTransactionSuccessRow, error)

	GetYearlyAmountSuccess(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyAmountTransactionSuccessRow, error)

	GetMonthlyAmountFailed(
		ctx context.Context,
		req *requests.MonthAmountTransaction,
	) ([]*db.GetMonthlyAmountTransactionFailedRow, error)

	GetYearlyAmountFailed(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyAmountTransactionFailedRow, error)

	GetMonthlyTransactionMethodSuccess(
		ctx context.Context,
		req *requests.MonthMethodTransaction,
	) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error)

	GetYearlyTransactionMethodSuccess(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTransactionMethodsSuccessRow, error)

	GetMonthlyTransactionMethodFailed(
		ctx context.Context,
		req *requests.MonthMethodTransaction,
	) ([]*db.GetMonthlyTransactionMethodsFailedRow, error)

	GetYearlyTransactionMethodFailed(
		ctx context.Context,
		year int,
	) ([]*db.GetYearlyTransactionMethodsFailedRow, error)
}

type TransactionStatsByMerchantRepository interface {
	GetMonthlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error)

	GetYearlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error)

	GetMonthlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error)

	GetYearlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error)

	GetMonthlyTransactionMethodByMerchantSuccess(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error)

	GetYearlyTransactionMethodByMerchantSuccess(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error)

	GetMonthlyTransactionMethodByMerchantFailed(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error)

	GetYearlyTransactionMethodByMerchantFailed(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error)
}

type TransactionQueryRepository interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsRow, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsActiveRow, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*db.GetTransactionsTrashedRow, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllTransactionByMerchant,
	) ([]*db.GetTransactionByMerchantRow, error)

	FindByID(
		ctx context.Context,
		transaction_id int,
	) (*db.GetTransactionByIDRow, error)

	FindByOrderID(
		ctx context.Context,
		order_id int,
	) (*db.GetTransactionByOrderIDRow, error)
}

type TransactionCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateTransactionRequest,
	) (*db.CreateTransactionRow, error)

	// CreateInTx runs the transaction insert inside the given database transaction
	// so the business row and its outbox event commit atomically.
	CreateInTx(
		ctx context.Context,
		tx pgx.Tx,
		request *requests.CreateTransactionRequest,
	) (*db.CreateTransactionRow, error)

	Update(
		ctx context.Context,
		request *requests.UpdateTransactionRequest,
	) (*db.UpdateTransactionRow, error)

	Trash(
		ctx context.Context,
		transaction_id int,
	) (*db.Transaction, error)

	Restore(
		ctx context.Context,
		transaction_id int,
	) (*db.Transaction, error)

	DeletePermanent(
		ctx context.Context,
		transaction_id int,
	) (bool, error)

	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
