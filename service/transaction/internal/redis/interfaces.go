package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type TransactionStatsCache interface {
	GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, bool)
	SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*response.TransactionMonthlyAmountSuccessResponse)

	GetCachedYearAmountSuccessCached(ctx context.Context, year int) ([]*response.TransactionYearlyAmountSuccessResponse, bool)
	SetCachedYearAmountSuccessCached(ctx context.Context, year int, res []*response.TransactionYearlyAmountSuccessResponse)

	GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, bool)
	SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*response.TransactionMonthlyAmountFailedResponse)

	GetCachedYearAmountFailedCached(ctx context.Context, year int) ([]*response.TransactionYearlyAmountFailedResponse, bool)
	SetCachedYearAmountFailedCached(ctx context.Context, year int, res []*response.TransactionYearlyAmountFailedResponse)

	GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, bool)
	SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*response.TransactionMonthlyMethodResponse)

	GetCachedYearMethodSuccessCached(ctx context.Context, year int) ([]*response.TransactionYearlyMethodResponse, bool)
	SetCachedYearMethodSuccessCached(ctx context.Context, year int, res []*response.TransactionYearlyMethodResponse)

	GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, bool)
	SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*response.TransactionMonthlyMethodResponse)

	GetCachedYearMethodFailedCached(ctx context.Context, year int) ([]*response.TransactionYearlyMethodResponse, bool)
	SetCachedYearMethodFailedCached(ctx context.Context, year int, res []*response.TransactionYearlyMethodResponse)
}

type TransactionStatsByMerchantCache interface {
	GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, bool)
	SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*response.TransactionMonthlyAmountSuccessResponse)

	GetCachedYearAmountSuccessCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, bool)
	SetCachedYearAmountSuccessCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*response.TransactionYearlyAmountSuccessResponse)

	GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, bool)
	SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*response.TransactionMonthlyAmountFailedResponse)

	GetCachedYearAmountFailedCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, bool)
	SetCachedYearAmountFailedCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*response.TransactionYearlyAmountFailedResponse)

	GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, bool)
	SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*response.TransactionMonthlyMethodResponse)

	GetCachedYearMethodSuccessCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, bool)
	SetCachedYearMethodSuccessCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*response.TransactionYearlyMethodResponse)

	GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, bool)
	SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*response.TransactionMonthlyMethodResponse)

	GetCachedYearMethodFailedCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, bool)
	SetCachedYearMethodFailedCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*response.TransactionYearlyMethodResponse)
}

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*response.TransactionResponse, total *int)

	GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*response.TransactionResponse, total *int)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*response.TransactionResponseDeleteAt, total *int)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*response.TransactionResponseDeleteAt, total *int)

	GetCachedTransactionCache(ctx context.Context, id int) (*response.TransactionResponse, bool)
	SetCachedTransactionCache(ctx context.Context, data *response.TransactionResponse)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*response.TransactionResponse, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *response.TransactionResponse)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
}
