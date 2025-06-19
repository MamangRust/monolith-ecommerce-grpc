package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type TransactionStatsCache interface {
	GetCachedMonthAmountSuccessCached(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, bool)
	SetCachedMonthAmountSuccessCached(req *requests.MonthAmountTransaction, res []*response.TransactionMonthlyAmountSuccessResponse)

	GetCachedYearAmountSuccessCached(year int) ([]*response.TransactionYearlyAmountSuccessResponse, bool)
	SetCachedYearAmountSuccessCached(year int, res []*response.TransactionYearlyAmountSuccessResponse)

	GetCachedMonthAmountFailedCached(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, bool)
	SetCachedMonthAmountFailedCached(req *requests.MonthAmountTransaction, res []*response.TransactionMonthlyAmountFailedResponse)

	GetCachedYearAmountFailedCached(year int) ([]*response.TransactionYearlyAmountFailedResponse, bool)
	SetCachedYearAmountFailedCached(year int, res []*response.TransactionYearlyAmountFailedResponse)

	GetCachedMonthMethodSuccessCached(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, bool)
	SetCachedMonthMethodSuccessCached(req *requests.MonthMethodTransaction, res []*response.TransactionMonthlyMethodResponse)

	GetCachedYearMethodSuccessCached(year int) ([]*response.TransactionYearlyMethodResponse, bool)
	SetCachedYearMethodSuccessCached(year int, res []*response.TransactionYearlyMethodResponse)

	GetCachedMonthMethodFailedCached(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, bool)
	SetCachedMonthMethodFailedCached(req *requests.MonthMethodTransaction, res []*response.TransactionMonthlyMethodResponse)

	GetCachedYearMethodFailedCached(year int) ([]*response.TransactionYearlyMethodResponse, bool)
	SetCachedYearMethodFailedCached(year int, res []*response.TransactionYearlyMethodResponse)
}

type TransactionStatsByMerchantCache interface {
	GetCachedMonthAmountSuccessCached(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, bool)
	SetCachedMonthAmountSuccessCached(req *requests.MonthAmountTransactionMerchant, res []*response.TransactionMonthlyAmountSuccessResponse)

	GetCachedYearAmountSuccessCached(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, bool)
	SetCachedYearAmountSuccessCached(req *requests.YearAmountTransactionMerchant, res []*response.TransactionYearlyAmountSuccessResponse)

	GetCachedMonthAmountFailedCached(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, bool)
	SetCachedMonthAmountFailedCached(req *requests.MonthAmountTransactionMerchant, res []*response.TransactionMonthlyAmountFailedResponse)

	GetCachedYearAmountFailedCached(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, bool)
	SetCachedYearAmountFailedCached(req *requests.YearAmountTransactionMerchant, res []*response.TransactionYearlyAmountFailedResponse)

	GetCachedMonthMethodSuccessCached(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, bool)
	SetCachedMonthMethodSuccessCached(req *requests.MonthMethodTransactionMerchant, res []*response.TransactionMonthlyMethodResponse)

	GetCachedYearMethodSuccessCached(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, bool)
	SetCachedYearMethodSuccessCached(req *requests.YearMethodTransactionMerchant, res []*response.TransactionYearlyMethodResponse)

	GetCachedMonthMethodFailedCached(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, bool)
	SetCachedMonthMethodFailedCached(req *requests.MonthMethodTransactionMerchant, res []*response.TransactionMonthlyMethodResponse)

	GetCachedYearMethodFailedCached(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, bool)
	SetCachedYearMethodFailedCached(req *requests.YearMethodTransactionMerchant, res []*response.TransactionYearlyMethodResponse)
}

type TransactionQueryCache interface {
	GetCachedTransactionsCache(req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, bool)
	SetCachedTransactionsCache(req *requests.FindAllTransaction, data []*response.TransactionResponse, total *int)

	GetCachedTransactionByMerchant(req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, bool)
	SetCachedTransactionByMerchant(req *requests.FindAllTransactionByMerchant, data []*response.TransactionResponse, total *int)

	GetCachedTransactionActiveCache(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, bool)
	SetCachedTransactionActiveCache(req *requests.FindAllTransaction, data []*response.TransactionResponseDeleteAt, total *int)

	GetCachedTransactionTrashedCache(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, bool)
	SetCachedTransactionTrashedCache(req *requests.FindAllTransaction, data []*response.TransactionResponseDeleteAt, total *int)

	GetCachedTransactionCache(id int) (*response.TransactionResponse, bool)
	SetCachedTransactionCache(data *response.TransactionResponse)

	GetCachedTransactionByOrderId(orderID int) (*response.TransactionResponse, bool)
	SetCachedTransactionByOrderId(orderID int, data *response.TransactionResponse)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(transactionID int)
}
