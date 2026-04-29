package merchant_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	merchantAllCacheKey      = "merchant:all:page:%d:pageSize:%d:search:%s"
	merchantByIdCacheKey     = "merchant:id:%d"
	merchantActiveCacheKey   = "merchant:active:page:%d:pageSize:%d:search:%s"
	merchantTrashedCacheKey  = "merchant:trashed:page:%d:pageSize:%d:search:%s"
	merchantByApiKeyCacheKey = "merchant:api_key:%s"
	merchantByUserIdCacheKey = "merchant:user_id:%d"

	ttlDefault = 5 * time.Minute
)

// Struktur pembungkus (merchantCachedResponseDB, dll.) sudah tidak diperlukan lagi
// karena tipe ApiResponse... sudah mencakup data dan paginasi.

type merchantQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantQueryCache(store *cache.CacheStore) *merchantQueryCache {
	return &merchantQueryCache{store: store}
}

func (m *merchantQueryCache) GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchant, bool) {
	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchant](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchant) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDeleteAt, bool) {
	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDeleteAt, bool) {
	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchant(ctx context.Context, id int) (*response.ApiResponseMerchant, bool) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	result, found := cache.GetFromCache[response.ApiResponseMerchant](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchant(ctx context.Context, data *response.ApiResponseMerchant) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByIdCacheKey, data.Data.ID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantsByUserId(ctx context.Context, id int) (*response.ApiResponsePaginationMerchant, bool) {
	key := fmt.Sprintf(merchantByUserIdCacheKey, id)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchant](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantQueryCache) SetCachedMerchantsByUserId(ctx context.Context, userId int, data *response.ApiResponsePaginationMerchant) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByUserIdCacheKey, userId)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
