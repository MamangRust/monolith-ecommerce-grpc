package order_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	orderAllCacheKey     = "order:all:page:%d:pageSize:%d:search:%s"
	orderByIdCacheKey    = "order:id:%d"
	orderActiveCacheKey  = "order:active:page:%d:pageSize:%d:search:%s"
	orderTrashedCacheKey = "order:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

// Struktur pembungkus (orderCacheResponseDB, dll.) sudah tidak diperlukan lagi
// karena tipe ApiResponse... sudah mencakup data dan paginasi.

type orderQueryCache struct {
	store *cache.CacheStore
}

func NewOrderQueryCache(store *cache.CacheStore) *orderQueryCache {
	return &orderQueryCache{store: store}
}

func (s *orderQueryCache) GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrder, bool) {
	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationOrder](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrder) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderQueryCache) GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrderDeleteAt, bool) {
	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationOrderDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrderDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderQueryCache) GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) (*response.ApiResponsePaginationOrderDeleteAt, bool) {
	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationOrderDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data *response.ApiResponsePaginationOrderDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderCache(ctx context.Context, order_id int) (*response.ApiResponseOrder, bool) {
	key := fmt.Sprintf(orderByIdCacheKey, order_id)
	result, found := cache.GetFromCache[response.ApiResponseOrder](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetCachedOrderCache(ctx context.Context, data *response.ApiResponseOrder) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderByIdCacheKey, data.Data.ID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
