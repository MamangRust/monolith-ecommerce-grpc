package merchantawards_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	merchantAwardAllCacheKey     = "merchant_award:all:page:%d:pageSize:%d:search:%s"
	merchantAwardByIdCacheKey    = "merchant_award:id:%d"
	merchantAwardActiveCacheKey  = "merchant_award:active:page:%d:pageSize:%d:search:%s"
	merchantAwardTrashedCacheKey = "merchant_award:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

// Struktur pembungkus (merchantAwardCacheResponseDB, dll.) sudah tidak diperlukan lagi
// karena tipe ApiResponse... sudah mencakup data dan paginasi.

type merchantAwardQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantAwardQueryCache(store *cache.CacheStore) *merchantAwardQueryCache {
	return &merchantAwardQueryCache{
		store: store,
	}
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantAward, bool) {
	key := fmt.Sprintf(merchantAwardAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantAward](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantAward) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantAwardAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantAwardDeleteAt, bool) {
	key := fmt.Sprintf(merchantAwardActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantAwardDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantAwardDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantAwardActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantAwardDeleteAt, bool) {
	key := fmt.Sprintf(merchantAwardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantAwardDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantAwardDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantAwardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAward(ctx context.Context, id int) (*response.ApiResponseMerchantAward, bool) {
	key := fmt.Sprintf(merchantAwardByIdCacheKey, id)

	result, found := cache.GetFromCache[response.ApiResponseMerchantAward](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAward(ctx context.Context, data *response.ApiResponseMerchantAward) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantAwardByIdCacheKey, data.Data.ID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
