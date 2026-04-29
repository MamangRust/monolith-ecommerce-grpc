package transaction_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type TransactionStatsCache interface {
	GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) (*response.ApiResponsesTransactionMonthSuccess, bool)
	SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, data *response.ApiResponsesTransactionMonthSuccess)

	GetCachedYearAmountSuccessCached(ctx context.Context, year int) (*response.ApiResponsesTransactionYearSuccess, bool)
	SetCachedYearAmountSuccessCached(ctx context.Context, year int, data *response.ApiResponsesTransactionYearSuccess)

	GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) (*response.ApiResponsesTransactionMonthFailed, bool)
	SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, data *response.ApiResponsesTransactionMonthFailed)

	GetCachedYearAmountFailedCached(ctx context.Context, year int) (*response.ApiResponsesTransactionYearFailed, bool)
	SetCachedYearAmountFailedCached(ctx context.Context, year int, data *response.ApiResponsesTransactionYearFailed)

	GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) (*response.ApiResponsesTransactionMonthMethod, bool)
	SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, data *response.ApiResponsesTransactionMonthMethod)

	GetCachedYearMethodSuccessCached(ctx context.Context, year int) (*response.ApiResponsesTransactionYearMethod, bool)
	SetCachedYearMethodSuccessCached(ctx context.Context, year int, data *response.ApiResponsesTransactionYearMethod)

	GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) (*response.ApiResponsesTransactionMonthMethod, bool)
	SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, data *response.ApiResponsesTransactionMonthMethod)

	GetCachedYearMethodFailedCached(ctx context.Context, year int) (*response.ApiResponsesTransactionYearMethod, bool)
	SetCachedYearMethodFailedCached(ctx context.Context, year int, data *response.ApiResponsesTransactionYearMethod)
}

type TransactionStatsByMerchantCache interface {
	GetCachedMonthAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) (*response.ApiResponsesTransactionMonthSuccess, bool)
	SetCachedMonthAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, data *response.ApiResponsesTransactionMonthSuccess)

	GetCachedYearAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) (*response.ApiResponsesTransactionYearSuccess, bool)
	SetCachedYearAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, data *response.ApiResponsesTransactionYearSuccess)

	GetCachedMonthAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) (*response.ApiResponsesTransactionMonthFailed, bool)
	SetCachedMonthAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, data *response.ApiResponsesTransactionMonthFailed)

	GetCachedYearAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) (*response.ApiResponsesTransactionYearFailed, bool)
	SetCachedYearAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, data *response.ApiResponsesTransactionYearFailed)

	GetCachedMonthMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) (*response.ApiResponsesTransactionMonthMethod, bool)
	SetCachedMonthMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, data *response.ApiResponsesTransactionMonthMethod)

	GetCachedYearMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) (*response.ApiResponsesTransactionYearMethod, bool)
	SetCachedYearMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, data *response.ApiResponsesTransactionYearMethod)

	GetCachedMonthMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) (*response.ApiResponsesTransactionMonthMethod, bool)
	SetCachedMonthMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, data *response.ApiResponsesTransactionMonthMethod)

	GetCachedYearMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) (*response.ApiResponsesTransactionYearMethod, bool)
	SetCachedYearMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, data *response.ApiResponsesTransactionYearMethod)
}

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) (*response.ApiResponsePaginationTransaction, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data *response.ApiResponsePaginationTransaction)

	GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) (*response.ApiResponsePaginationTransaction, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data *response.ApiResponsePaginationTransaction)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) (*response.ApiResponsePaginationTransaction, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data *response.ApiResponsePaginationTransaction)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) (*response.ApiResponsePaginationTransactionDeleteAt, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data *response.ApiResponsePaginationTransactionDeleteAt)

	GetCachedTransactionCache(ctx context.Context, id int) (*response.ApiResponseTransaction, bool)
	SetCachedTransactionCache(ctx context.Context, data *response.ApiResponseTransaction)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*response.ApiResponseTransaction, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *response.ApiResponseTransaction)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
	InvalidateTransactionCache(ctx context.Context)
}
