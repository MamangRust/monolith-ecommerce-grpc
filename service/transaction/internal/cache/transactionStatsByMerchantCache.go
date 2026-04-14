package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

const (
	transactonMonthlyAmountSuccessByMerchantKey = "transaction:month:amount:success:merchant:%d:month:%d:year:%d"
	transactonMonthlyAmountFailedByMerchantKey  = "transaction:month:amount:failed:merchant:%d:month:%d:year:%d"
	transactonYearlyAmountSuccessByMerchantKey  = "transaction:year:amount:success:merchant:%d:year:%d"
	transactonYearlyAmountFailedByMerchantKey   = "transaction:year:amount:failed:merchant:%d:year:%d"
	transactonMonthlyMethodSuccessByMerchantKey = "transaction:month:method:success:merchant:%d:month:%d:year:%d"
	transactonMonthlyMethodFailedByMerchantKey  = "transaction:month:method:failed:merchant:%d:month:%d:year:%d"
	transactonYearlyMethodSuccessByMerchantKey  = "transaction:year:method:success:merchant:%d:year:%d"
	transactonYearlyMethodFailedByMerchantKey   = "transaction:year:method:failed:merchant:%d:year:%d"
)

type transactionStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByMerchantCache(store *cache.CacheStore) *transactionStatsByMerchantCache {
	return &transactionStatsByMerchantCache{store: store}
}

func (t *transactionStatsByMerchantCache) GetCachedMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthlyAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*db.GetMonthlyAmountTransactionSuccessByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthlyAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearlyAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*db.GetYearlyAmountTransactionSuccessByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearlyAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthlyAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyAmountTransactionFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*db.GetMonthlyAmountTransactionFailedByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthlyAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearlyAmountFailedByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyAmountTransactionFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*db.GetYearlyAmountTransactionFailedByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearlyAmountFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthlyMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, bool) {
	key := fmt.Sprintf(transactonMonthlyMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthlyMethodSuccessByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthlyMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearlyMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, bool) {
	key := fmt.Sprintf(transactonYearlyMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsByMerchantSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearlyMethodSuccessByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*db.GetYearlyTransactionMethodsByMerchantSuccessRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearlyMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthlyMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, bool) {
	key := fmt.Sprintf(transactonMonthlyMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTransactionMethodsByMerchantFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthlyMethodFailedByMerchant(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*db.GetMonthlyTransactionMethodsByMerchantFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonMonthlyMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearlyMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, bool) {
	key := fmt.Sprintf(transactonYearlyMethodFailedByMerchantKey, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTransactionMethodsByMerchantFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearlyMethodFailedByMerchant(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*db.GetYearlyTransactionMethodsByMerchantFailedRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(transactonYearlyMethodFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
