package mencache

import (
	"context"
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

func (t *transactionStatsCache) GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, bool) {
	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyAmountSuccessResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*response.TransactionMonthlyAmountSuccessResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountSuccessCached(ctx context.Context, year int) ([]*response.TransactionYearlyAmountSuccessResponse, bool) {
	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyAmountSuccessResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearAmountSuccessCached(ctx context.Context, year int, res []*response.TransactionYearlyAmountSuccessResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, bool) {
	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyAmountFailedResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*response.TransactionMonthlyAmountFailedResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountFailedCached(ctx context.Context, year int) ([]*response.TransactionYearlyAmountFailedResponse, bool) {
	key := fmt.Sprintf(transactionYearAmountFailedKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyAmountFailedResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearAmountFailedCached(ctx context.Context, year int, res []*response.TransactionYearlyAmountFailedResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountFailedKey, year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, bool) {
	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*response.TransactionMonthlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodSuccessCached(ctx context.Context, year int) ([]*response.TransactionYearlyMethodResponse, bool) {
	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearMethodSuccessCached(ctx context.Context, year int, res []*response.TransactionYearlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, bool) {
	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*response.TransactionMonthlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodFailedCached(ctx context.Context, year int) ([]*response.TransactionYearlyMethodResponse, bool) {
	key := fmt.Sprintf(transactionYearMethodFailedKey, year)

	result, found := GetFromCache[[]*response.TransactionYearlyMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearMethodFailedCached(ctx context.Context, year int, res []*response.TransactionYearlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodFailedKey, year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}
