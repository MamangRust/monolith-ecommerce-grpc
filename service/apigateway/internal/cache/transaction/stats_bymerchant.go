package transaction_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

const (
	transactonMonthAmountSuccessByMerchantKey = "transaction:month:amount:success:merchant:%d:month:%d:year:%d"
	transactonMonthAmountFailedByMerchantKey  = "transaction:month:amount:failed:merchant:%d:month:%d:year:%d"
	transactonYearAmountSuccessByMerchantKey  = "transaction:year:amount:success:merchant:%d:year:%d"
	transactonYearAmountFailedByMerchantKey   = "transaction:year:amount:failed:merchant:%d:year:%d"
	transactonMonthMethodSuccessByMerchantKey = "transaction:month:method:success:merchant:%d:month:%d:year:%d"
	transactonMonthMethodFailedByMerchantKey  = "transaction:month:method:failed:merchant:%d:month:%d:year:%d"
	transactonYearMethodSuccessByMerchantKey  = "transaction:year:method:success:merchant:%d:year:%d"
	transactonYearMethodFailedByMerchantKey   = "transaction:year:method:failed:merchant:%d:year:%d"
)

type transactionStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByMerchantCache(store *cache.CacheStore) *transactionStatsByMerchantCache {
	return &transactionStatsByMerchantCache{store: store}
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) (*response.ApiResponsesTransactionMonthSuccess, bool) {
	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionMonthSuccess](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, data *response.ApiResponsesTransactionMonthSuccess) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) (*response.ApiResponsesTransactionYearSuccess, bool) {
	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionYearSuccess](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, data *response.ApiResponsesTransactionYearSuccess) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) (*response.ApiResponsesTransactionMonthFailed, bool) {
	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionMonthFailed](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, data *response.ApiResponsesTransactionMonthFailed) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) (*response.ApiResponsesTransactionYearFailed, bool) {
	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionYearFailed](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, data *response.ApiResponsesTransactionYearFailed) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) (*response.ApiResponsesTransactionMonthMethod, bool) {
	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionMonthMethod](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, data *response.ApiResponsesTransactionMonthMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) (*response.ApiResponsesTransactionYearMethod, bool) {
	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionYearMethod](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, data *response.ApiResponsesTransactionYearMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) (*response.ApiResponsesTransactionMonthMethod, bool) {
	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionMonthMethod](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, data *response.ApiResponsesTransactionMonthMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) (*response.ApiResponsesTransactionYearMethod, bool) {
	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponsesTransactionYearMethod](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, data *response.ApiResponsesTransactionYearMethod) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}
