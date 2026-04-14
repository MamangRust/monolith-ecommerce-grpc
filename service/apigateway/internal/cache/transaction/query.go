package transaction_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	transactionAllCacheKey        = "transaction:all:page:%d:pageSize:%d:search:%s"
	transactionByIdCacheKey       = "transaction:id:%d"
	transactionByMerchantCacheKey = "transaction:merchant:%d:page:%d:pageSize:%d:search:%s"
	transactionActiveCacheKey     = "transaction:active:page:%d:pageSize:%d:search:%s"
	transactionTrashedCacheKey    = "transaction:trashed:page:%d:pageSize:%d:search:%s"
	transactionByOrderCacheKey    = "transaction:order:%d"
	ttlDefault                    = 5 * time.Minute
)

// Struktur pembungkus (transactionCacheResponseDB, dll.) sudah tidak diperlukan lagi
// karena tipe ApiResponse... sudah mencakup data dan paginasi.

type transactionQueryCache struct {
	store *cache.CacheStore
}

func NewTransactionQueryCache(store *cache.CacheStore) *transactionQueryCache {
	return &transactionQueryCache{store: store}
}

func (t *transactionQueryCache) GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) (*response.ApiResponsePaginationTransaction, bool) {
	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data *response.ApiResponsePaginationTransaction) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) (*response.ApiResponsePaginationTransaction, bool) {
	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data *response.ApiResponsePaginationTransaction) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) (*response.ApiResponsePaginationTransaction, bool) {
	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data *response.ApiResponsePaginationTransaction) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) (*response.ApiResponsePaginationTransactionDeleteAt, bool) {
	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationTransactionDeleteAt](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data *response.ApiResponsePaginationTransactionDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionCache(ctx context.Context, id int) (*response.ApiResponseTransaction, bool) {
	key := fmt.Sprintf(transactionByIdCacheKey, id)
	result, found := cache.GetFromCache[response.ApiResponseTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionCache(ctx context.Context, data *response.ApiResponseTransaction) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionByIdCacheKey, data.Data.ID)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*response.ApiResponseTransaction, bool) {
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)
	result, found := cache.GetFromCache[response.ApiResponseTransaction](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *response.ApiResponseTransaction) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)

	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}
