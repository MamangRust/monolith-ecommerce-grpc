package cache

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	shippingAddressAllCacheKey     = "shipping_address:all:page:%d:pageSize:%d:search:%s"
	shippingAddressActiveCacheKey  = "shipping_address:active:page:%d:pageSize:%d:search:%s"
	shippingAddressTrashedCacheKey = "shipping_address:trashed:page:%d:pageSize:%d:search:%s"

	shippingAddressByOrderIdCacheKey = "shipping_address:order_id:%d"
	shippingAddressByIdCacheKey      = "shipping_address:id:%d"

	ttlDefault = 5 * time.Minute
)

type shippingAddressCacheResponseDB struct {
	Data  []*db.GetShippingAddressRow `json:"data"`
	Total *int                        `json:"total_records"`
}

type shippingAddressActiveCacheResponseDB struct {
	Data  []*db.GetShippingAddressActiveRow `json:"data"`
	Total *int                              `json:"total_records"`
}

type shippingAddressTrashedCacheResponseDB struct {
	Data  []*db.GetShippingAddressTrashedRow `json:"data"`
	Total *int                               `json:"total_records"`
}

type shippingAddressQueryCache struct {
	store *cache.CacheStore
}

func NewShippingAddressQueryCache(store *cache.CacheStore) *shippingAddressQueryCache {
	return &shippingAddressQueryCache{store: store}
}

func (r *shippingAddressQueryCache) GetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressRow, *int, bool) {
	key := fmt.Sprintf(shippingAddressAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[shippingAddressCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *shippingAddressQueryCache) SetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*db.GetShippingAddressRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if res == nil {
		res = []*db.GetShippingAddressRow{}
	}

	key := fmt.Sprintf(shippingAddressAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &shippingAddressCacheResponseDB{Data: res, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *shippingAddressQueryCache) GetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressTrashedRow, *int, bool) {
	key := fmt.Sprintf(shippingAddressTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[shippingAddressTrashedCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *shippingAddressQueryCache) SetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*db.GetShippingAddressTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if res == nil {
		res = []*db.GetShippingAddressTrashedRow{}
	}

	key := fmt.Sprintf(shippingAddressTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &shippingAddressTrashedCacheResponseDB{Data: res, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *shippingAddressQueryCache) GetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressActiveRow, *int, bool) {
	key := fmt.Sprintf(shippingAddressActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[shippingAddressActiveCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *shippingAddressQueryCache) SetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*db.GetShippingAddressActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if res == nil {
		res = []*db.GetShippingAddressActiveRow{}
	}

	key := fmt.Sprintf(shippingAddressActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &shippingAddressActiveCacheResponseDB{Data: res, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *shippingAddressQueryCache) GetCachedShippingAddressCache(ctx context.Context, shipping_id int) (*db.GetShippingByIDRow, bool) {
	key := fmt.Sprintf(shippingAddressByIdCacheKey, shipping_id)
	result, found := cache.GetFromCache[db.GetShippingByIDRow](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *shippingAddressQueryCache) SetCachedShippingAddressCache(ctx context.Context, data *db.GetShippingByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shippingAddressByIdCacheKey, data.ShippingAddressID)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *shippingAddressQueryCache) GetCachedShippingAddressByOrderCache(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, bool) {
	key := fmt.Sprintf(shippingAddressByOrderIdCacheKey, order_id)
	result, found := cache.GetFromCache[db.GetShippingAddressByOrderIDRow](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *shippingAddressQueryCache) SetCachedShippingAddressByOrderCache(ctx context.Context, data *db.GetShippingAddressByOrderIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shippingAddressByOrderIdCacheKey, data.OrderID) // Typo diperbaiki
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}
