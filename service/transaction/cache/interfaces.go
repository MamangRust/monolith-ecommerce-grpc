package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type TransactionStatsCache interface {
	GetCachedMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, bool)
	SetCachedMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction, res []*db.GetMonthlyAmountTransactionSuccessRow)

	GetCachedYearlyAmountSuccess(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, bool)
	SetCachedYearlyAmountSuccess(ctx context.Context, year int, res []*db.GetYearlyAmountTransactionSuccessRow)

	GetCachedMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, bool)
	SetCachedMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction, res []*db.GetMonthlyAmountTransactionFailedRow)

	GetCachedYearlyAmountFailed(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, bool)
	SetCachedYearlyAmountFailed(ctx context.Context, year int, res []*db.GetYearlyAmountTransactionFailedRow)

	GetCachedMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, bool)
	SetCachedMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction, res []*db.GetMonthlyTransactionMethodsSuccessRow)

	GetCachedYearlyMethodSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, bool)
	SetCachedYearlyMethodSuccess(ctx context.Context, year int, res []*db.GetYearlyTransactionMethodsSuccessRow)

	GetCachedMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, bool)
	SetCachedMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction, res []*db.GetMonthlyTransactionMethodsFailedRow)

	GetCachedYearlyMethodFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, bool)
	SetCachedYearlyMethodFailed(ctx context.Context, year int, res []*db.GetYearlyTransactionMethodsFailedRow)
}

type TransactionStatsByMerchantCache interface {
	GetCachedMonthlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, bool)

	SetCachedMonthlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
		res []*db.GetMonthlyAmountTransactionSuccessByMerchantRow,
	)

	GetCachedYearlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, bool)

	SetCachedYearlyAmountSuccessByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
		res []*db.GetYearlyAmountTransactionSuccessByMerchantRow,
	)

	GetCachedMonthlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
	) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, bool)

	SetCachedMonthlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.MonthAmountTransactionMerchant,
		res []*db.GetMonthlyAmountTransactionFailedByMerchantRow,
	)

	GetCachedYearlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
	) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, bool)

	SetCachedYearlyAmountFailedByMerchant(
		ctx context.Context,
		req *requests.YearAmountTransactionMerchant,
		res []*db.GetYearlyAmountTransactionFailedByMerchantRow,
	)

	GetCachedMonthlyMethodSuccessByMerchant(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, bool)

	SetCachedMonthlyMethodSuccessByMerchant(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
		res []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow,
	)

	GetCachedYearlyMethodSuccessByMerchant(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, bool)

	SetCachedYearlyMethodSuccessByMerchant(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
		res []*db.GetYearlyTransactionMethodsByMerchantSuccessRow,
	)

	GetCachedMonthlyMethodFailedByMerchant(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
	) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, bool)

	SetCachedMonthlyMethodFailedByMerchant(
		ctx context.Context,
		req *requests.MonthMethodTransactionMerchant,
		res []*db.GetMonthlyTransactionMethodsByMerchantFailedRow,
	)

	GetCachedYearlyMethodFailedByMerchant(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
	) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, bool)

	SetCachedYearlyMethodFailedByMerchant(
		ctx context.Context,
		req *requests.YearMethodTransactionMerchant,
		res []*db.GetYearlyTransactionMethodsByMerchantFailedRow,
	)
}

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsRow, total *int)

	GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*db.GetTransactionByMerchantRow, *int, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*db.GetTransactionByMerchantRow, total *int)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsActiveRow, total *int)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsTrashedRow, total *int)

	GetCachedTransactionCache(ctx context.Context, id int) (*db.GetTransactionByIDRow, bool)
	SetCachedTransactionCache(ctx context.Context, data *db.GetTransactionByIDRow)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*db.GetTransactionByOrderIDRow, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *db.GetTransactionByOrderIDRow)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
	InvalidateTransactionCache(ctx context.Context)
}
