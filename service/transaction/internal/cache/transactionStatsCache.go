package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

const (
	transactionMonthlyAmountSuccessKey = "transaction:month:amount:success:month:%d:year:%d"
	transactionMonthlyAmountFailedKey  = "transaction:month:amount:failed:month:%d:year:%d"
	transactionYearlyAmountSuccessKey  = "transaction:year:amount:success:year:%d"
	transactionYearlyAmountFailedKey   = "transaction:year:amount:failed:year:%d"
	transactionMonthlyMethodSuccessKey = "transaction:month:method:success:month:%d:year:%d"
	transactionMonthlyMethodFailedKey  = "transaction:month:method:failed:month:%d:year:%d"
	transactionYearlyMethodSuccessKey  = "transaction:year:method:success:year:%d"
	transactionYearlyMethodFailedKey   = "transaction:year:method:failed:year:%d"
)

type transactionStatsCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsCache(store *cache.CacheStore) *transactionStatsCache {
	return &transactionStatsCache{store: store}
}

func (t *transactionStatsCache) GetCachedMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, bool) {
	key := fmt.Sprintf(transactionMonthlyAmountSuccessKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction, res []*db.GetMonthlyAmountTransactionSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthlyAmountSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearlyAmountSuccess(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, bool) {
	key := fmt.Sprintf(transactionYearlyAmountSuccessKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedYearlyAmountSuccess(ctx context.Context, year int, res []*db.GetYearlyAmountTransactionSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionYearlyAmountSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, bool) {
	key := fmt.Sprintf(transactionMonthlyAmountFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction, res []*db.GetMonthlyAmountTransactionFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthlyAmountFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearlyAmountFailed(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, bool) {
	key := fmt.Sprintf(transactionYearlyAmountFailedKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedYearlyAmountFailed(ctx context.Context, year int, res []*db.GetYearlyAmountTransactionFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionYearlyAmountFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, bool) {
	key := fmt.Sprintf(transactionMonthlyMethodSuccessKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction, res []*db.GetMonthlyTransactionMethodsSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthlyMethodSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearlyMethodSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, bool) {
	key := fmt.Sprintf(transactionYearlyMethodSuccessKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedYearlyMethodSuccess(ctx context.Context, year int, res []*db.GetYearlyTransactionMethodsSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionYearlyMethodSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, bool) {
	key := fmt.Sprintf(transactionMonthlyMethodFailedKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction, res []*db.GetMonthlyTransactionMethodsFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionMonthlyMethodFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearlyMethodFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, bool) {
	key := fmt.Sprintf(transactionYearlyMethodFailedKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsCache) SetCachedYearlyMethodFailed(ctx context.Context, year int, res []*db.GetYearlyTransactionMethodsFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactionYearlyMethodFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
