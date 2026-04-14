package shippingaddress_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	shippingAddressAllCacheKey     = "shipping_address:all:page:%d:pageSize:%d:search:%s"
	shippingAddressActiveCacheKey  = "shipping_address:active:page:%d:pageSize:%d:search:%s"
	shippingAddressTrashedCacheKey = "shipping_address:trashed:page:%d:pageSize:%d:search:%s"

	shippingAddressByOrderIdCacheKey = "shipping_address:order_id:%d"
	shippingAddressByIdCacheKey      = "shipping_address:id:%d"

	ttlDefault = 5 * time.Minute
)

type shippingAddressQueryCache struct {
	store *cache.CacheStore
}

func NewShippingAddressQueryCache(store *cache.CacheStore) *shippingAddressQueryCache {
	return &shippingAddressQueryCache{store: store}
}

func (r *shippingAddressQueryCache) GetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress) (*response.ApiResponsePaginationShippingAddress, bool) {
	key := fmt.Sprintf(shippingAddressAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationShippingAddress](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *shippingAddressQueryCache) SetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress, data *response.ApiResponsePaginationShippingAddress) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shippingAddressAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *shippingAddressQueryCache) GetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress) (*response.ApiResponsePaginationShippingAddressDeleteAt, bool) {
	key := fmt.Sprintf(shippingAddressActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationShippingAddressDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *shippingAddressQueryCache) SetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress, data *response.ApiResponsePaginationShippingAddressDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shippingAddressActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *shippingAddressQueryCache) GetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress) (*response.ApiResponsePaginationShippingAddressDeleteAt, bool) {
	key := fmt.Sprintf(shippingAddressTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationShippingAddressDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *shippingAddressQueryCache) SetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress, data *response.ApiResponsePaginationShippingAddressDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shippingAddressTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *shippingAddressQueryCache) GetCachedShippingAddressCache(ctx context.Context, shipping_id int) (*response.ApiResponseShippingAddress, bool) {
	key := fmt.Sprintf(shippingAddressByIdCacheKey, shipping_id)
	result, found := cache.GetFromCache[response.ApiResponseShippingAddress](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *shippingAddressQueryCache) SetCachedShippingAddressCache(ctx context.Context, data *response.ApiResponseShippingAddress) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shippingAddressByIdCacheKey, data.Data.ID)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *shippingAddressQueryCache) GetCachedShippingAddressByOrderCache(ctx context.Context, order_id int) (*response.ApiResponseShippingAddress, bool) {
	key := fmt.Sprintf(shippingAddressByOrderIdCacheKey, order_id)
	result, found := cache.GetFromCache[response.ApiResponseShippingAddress](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *shippingAddressQueryCache) SetCachedShippingAddressByOrderCache(ctx context.Context, data *response.ApiResponseShippingAddress) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shippingAddressByOrderIdCacheKey, data.Data.OrderID)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}
