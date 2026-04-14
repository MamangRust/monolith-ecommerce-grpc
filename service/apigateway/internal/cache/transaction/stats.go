package transaction_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

const (
	transactionMonthAmountSuccessKey = "transaction:month:amount:success:month:%d:year:%d"
	transactionMonthAmountFailedKey  = "transaction:month:amount:failed:month:%d:year:%d"
	transactionYearAmountSuccessKey  = "transaction:year:amount:success:year:%d"
	transactionYearAmountFailedKey   = "transaction:year:amount:failed:year:%d"
	transactionMonthMethodSuccessKey = "transaction:month:method:success:month:%d:year:%d"
	transactionMonthMethodFailedKey  = "transaction:month:method:failed:month:%d:year:%d"
	transactionYearMethodSuccessKey  = "transaction:year:method:success:year:%d"
	transactionYearMethodFailedKey   = "transaction:year:method:failed:year:%d"
)

type transactionStatsCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsCache(store *cache.CacheStore) *transactionStatsCache {
	return &transactionStatsCache{store: store}
}

func (t *transactionStatsCache) GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) (*response.ApiResponsesTransactionMonthSuccess, bool) {
	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionMonthSuccess](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, data *response.ApiResponsesTransactionMonthSuccess) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountSuccessCached(ctx context.Context, year int) (*response.ApiResponsesTransactionYearSuccess, bool) {
	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionYearSuccess](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsCache) SetCachedYearAmountSuccessCached(ctx context.Context, year int, data *response.ApiResponsesTransactionYearSuccess) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) (*response.ApiResponsesTransactionMonthFailed, bool) {
	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionMonthFailed](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, data *response.ApiResponsesTransactionMonthFailed) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountFailedCached(ctx context.Context, year int) (*response.ApiResponsesTransactionYearFailed, bool) {
	key := fmt.Sprintf(transactionYearAmountFailedKey, year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionYearFailed](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsCache) SetCachedYearAmountFailedCached(ctx context.Context, year int, data *response.ApiResponsesTransactionYearFailed) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionYearAmountFailedKey, year)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) (*response.ApiResponsesTransactionMonthMethod, bool) {
	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionMonthMethod](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, data *response.ApiResponsesTransactionMonthMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodSuccessCached(ctx context.Context, year int) (*response.ApiResponsesTransactionYearMethod, bool) {
	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionYearMethod](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsCache) SetCachedYearMethodSuccessCached(ctx context.Context, year int, data *response.ApiResponsesTransactionYearMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) (*response.ApiResponsesTransactionMonthMethod, bool) {
	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionMonthMethod](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, data *response.ApiResponsesTransactionMonthMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodFailedCached(ctx context.Context, year int) (*response.ApiResponsesTransactionYearMethod, bool) {
	key := fmt.Sprintf(transactionYearMethodFailedKey, year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionYearMethod](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsCache) SetCachedYearMethodFailedCached(ctx context.Context, year int, data *response.ApiResponsesTransactionYearMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionYearMethodFailedKey, year)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}
