package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	shippingAddressAllCacheKey     = "shipping_address:all:page:%d:pageSize:%d:search:%s"
	shippingAddressActiveCacheKey  = "shipping_address:active:page:%d:pageSize:%d:search:%s"
	shippingAddressTrashedCacheKey = "shipping_address:trashed:page:%d:pageSize:%d:search:%s"

	shiippingAddressByOrderIdCacheKey = "shipping_address:order_id:%d"
	shippingAddressByIdCacheKey       = "shipping_address:id:%d"

	ttlDefault = 5 * time.Minute
)

type shippingAddressCacheResponse struct {
	Data  []*response.ShippingAddressResponse `json:"data"`
	Total *int                                `json:"total_records"`
}

type shippingAddressCacheResponseDeleteAt struct {
	Data  []*response.ShippingAddressResponseDeleteAt `json:"data"`
	Total *int                                        `json:"total_records"`
}

type shippingAddressQueryCache struct {
	store *CacheStore
}

func NewShippingAddressQueryCache(store *CacheStore) *shippingAddressQueryCache {
	return &shippingAddressQueryCache{store: store}
}

func (r *shippingAddressQueryCache) GetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, bool) {
	key := fmt.Sprintf(shippingAddressAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[shippingAddressCacheResponse](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *shippingAddressQueryCache) SetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if res == nil {
		res = []*response.ShippingAddressResponse{}
	}

	key := fmt.Sprintf(shippingAddressAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &shippingAddressCacheResponse{Data: res, Total: total}
	SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *shippingAddressQueryCache) GetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(shippingAddressTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[shippingAddressCacheResponseDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *shippingAddressQueryCache) SetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if res == nil {
		res = []*response.ShippingAddressResponseDeleteAt{}
	}

	key := fmt.Sprintf(shippingAddressTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &shippingAddressCacheResponseDeleteAt{Data: res, Total: total}
	SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *shippingAddressQueryCache) GetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(shippingAddressActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[shippingAddressCacheResponseDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *shippingAddressQueryCache) SetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if res == nil {
		res = []*response.ShippingAddressResponseDeleteAt{}
	}

	key := fmt.Sprintf(shippingAddressActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &shippingAddressCacheResponseDeleteAt{Data: res, Total: total}

	SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *shippingAddressQueryCache) GetCachedShippingAddressCache(ctx context.Context, shipping_id int) (*response.ShippingAddressResponse, bool) {
	key := fmt.Sprintf(shippingAddressByIdCacheKey, shipping_id)
	result, found := GetFromCache[*response.ShippingAddressResponse](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *shippingAddressQueryCache) SetCachedShippingAddressCache(ctx context.Context, data *response.ShippingAddressResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shippingAddressByIdCacheKey, data.ID)
	SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *shippingAddressQueryCache) GetCachedShippingAddressByOrderCache(ctx context.Context, order_id int) (*response.ShippingAddressResponse, bool) {
	key := fmt.Sprintf(shiippingAddressByOrderIdCacheKey, order_id)
	result, found := GetFromCache[*response.ShippingAddressResponse](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *shippingAddressQueryCache) SetCachedShippingAddressByOrderCache(ctx context.Context, data *response.ShippingAddressResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(shiippingAddressByOrderIdCacheKey, data.OrderID)
	SetToCache(ctx, r.store, key, data, ttlDefault)
}
