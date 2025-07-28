package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	transactonMonthAmountSuccessByMerchantKey = "transaction:month:amount:success:merchant:%d:month:%d:year:%d"
	transactonMonthAmountFailedByMerchantKey  = "transaction:month:amount:failed:merchant:%d:month:%d:year:%d"

	transactonYearAmountSuccessByMerchantKey = "transaction:year:amount:success:merchant:%d:year:%d"
	transactonYearAmountFailedByMerchantKey  = "transaction:year:amount:failed:merchant:%d:year:%d"

	transactonMonthMethodSuccessByMerchantKey = "transaction:month:method:success:merchant:%d:month:%d:year:%d"
	transactonMonthMethodFailedByMerchantKey  = "transaction:month:method:failed:merchant:%d:month:%d:year:%d"

	transactonYearMethodSuccessByMerchantKey = "transaction:year:method:success:merchant:%d:year:%d"
	transactonYearMethodFailedByMerchantKey  = "transaction:year:method:failed:merchant:%d:year:%d"
)

type transactionStatsByMerchantCache struct {
	store *CacheStore
}

func NewTransactionStatsByMerchantCache(store *CacheStore) *transactionStatsByMerchantCache {
	return &transactionStatsByMerchantCache{store: store}
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, bool) {
	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyAmountSuccessResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*response.TransactionMonthlyAmountSuccessResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, bool) {
	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyAmountFailedResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*response.TransactionMonthlyAmountFailedResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountFailedCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, bool) {
	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.TransactionYearlyAmountFailedResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountFailedCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*response.TransactionYearlyAmountFailedResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountSuccessCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, bool) {
	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.TransactionYearlyAmountSuccessResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountSuccessCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*response.TransactionYearlyAmountSuccessResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, bool) {
	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*response.TransactionMonthlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodSuccessCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, bool) {
	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.TransactionYearlyMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodSuccessCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*response.TransactionYearlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, bool) {
	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.TransactionMonthlyMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*response.TransactionMonthlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodFailedCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, bool) {
	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.TransactionYearlyMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodFailedCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*response.TransactionYearlyMethodResponse) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)

	SetToCache(ctx, t.store, key, &res, ttlDefault)
}
