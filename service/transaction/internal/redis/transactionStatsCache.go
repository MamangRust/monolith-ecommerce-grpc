package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	transactionMonthAmountSuccessKey = "transaction:month:amount:success:month:%d:year:%d"
	transactionMonthAmountFailedKey  = "transaction:month:amount:failed:month:%d:year:%d"

	transactionYearAmountSuccessKey = "transaction:year:amount:success:year:%d"
	transactionYearAmountFailedKey  = "transaction:year:amount:failed:year:%d"

	transactionMonthMethodSuccessKey = "transaction:month:method:success:month:%d:year:%d"
	transactionMonthMethodFailedKey  = "transaction:month:method:failed:month:%d:year:%d"

	transactionYearMethodSuccessKey = "transaction:year:method:success:year:%d"
	transactionYearMethodFailedKey  = "transaction:year:method:failed:year:%d"
)

type transactionStatsCache struct {
	store *CacheStore
}

func NewTransactionStatsCache(store *CacheStore) *transactionStatsCache {
	return &transactionStatsCache{store: store}
}

func (t *transactionStatsCache) GetCachedMonthAmountSuccessCached(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, bool) {
	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyAmountSuccessResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountSuccessCached(req *requests.MonthAmountTransaction, res []*response.TransactionMonthlyAmountSuccessResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	SetToCache(t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountSuccessCached(year int) ([]*response.TransactionYearlyAmountSuccessResponse, bool) {
	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyAmountSuccessResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearAmountSuccessCached(year int, res []*response.TransactionYearlyAmountSuccessResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)

	SetToCache(t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthAmountFailedCached(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, bool) {
	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyAmountFailedResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountFailedCached(req *requests.MonthAmountTransaction, res []*response.TransactionMonthlyAmountFailedResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)

	SetToCache(t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountFailedCached(year int) ([]*response.TransactionYearlyAmountFailedResponse, bool) {
	key := fmt.Sprintf(transactionYearAmountFailedKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyAmountFailedResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearAmountFailedCached(year int, res []*response.TransactionYearlyAmountFailedResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountFailedKey, year)

	SetToCache(t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodSuccessCached(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, bool) {
	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyMethodResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodSuccessCached(req *requests.MonthMethodTransaction, res []*response.TransactionMonthlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)

	SetToCache(t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodSuccessCached(year int) ([]*response.TransactionYearlyMethodResponse, bool) {
	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyMethodResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearMethodSuccessCached(year int, res []*response.TransactionYearlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)

	SetToCache(t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodFailedCached(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, bool) {
	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyMethodResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodFailedCached(req *requests.MonthMethodTransaction, res []*response.TransactionMonthlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)

	SetToCache(t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodFailedCached(year int) ([]*response.TransactionYearlyMethodResponse, bool) {
	key := fmt.Sprintf(transactionYearMethodFailedKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyMethodResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearMethodFailedCached(year int, res []*response.TransactionYearlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodFailedKey, year)

	SetToCache(t.store, key, &res, ttlDefault)
}
