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
	orderAllCacheKey     = "order:all:page:%d:pageSize:%d:search:%s"
	orderByIdCacheKey    = "order:id:%d"
	orderActiveCacheKey  = "order:active:page:%d:pageSize:%d:search:%s"
	orderTrashedCacheKey = "order:trashed:page:%d:pageSize:%d:search:%s"
	orderByMerchantCacheKey = "order:merchant:merchantID:%d:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type orderCacheResponseDB struct {
	Data         []*db.GetOrdersRow `json:"data"`
	TotalRecords *int               `json:"total_records"`
}

type orderActiveCacheResponseDB struct {
	Data         []*db.GetOrdersActiveRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type orderTrashedCacheResponseDB struct {
	Data         []*db.GetOrdersTrashedRow `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type orderMerchantCacheResponseDB struct {
	Data         []*db.GetOrdersByMerchantRow `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type orderQueryCache struct {
	store *cache.CacheStore
}

func NewOrderQueryCache(store *cache.CacheStore) *orderQueryCache {
	return &orderQueryCache{store: store}
}

func (s *orderQueryCache) GetOrderAllCache(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersRow, *int, bool) {
	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[orderCacheResponseDB](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderAllCache(ctx context.Context, req *requests.FindAllOrder, data []*db.GetOrdersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetOrdersRow{}
	}

	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersActiveRow, *int, bool) {
	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[orderActiveCacheResponseDB](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrder, data []*db.GetOrdersActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetOrdersActiveRow{}
	}

	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersTrashedRow, *int, bool) {
	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[orderTrashedCacheResponseDB](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrder, data []*db.GetOrdersTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetOrdersTrashedRow{}
	}

	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderCache(ctx context.Context, order_id int) (*db.GetOrderByIDRow, bool) {
	key := fmt.Sprintf(orderByIdCacheKey, order_id)
	result, found := cache.GetFromCache[db.GetOrderByIDRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderQueryCache) SetCachedOrderCache(ctx context.Context, data *db.GetOrderByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderByIdCacheKey, data.OrderID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderQueryCache) GetOrderByMerchantCache(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*db.GetOrdersByMerchantRow, *int, bool) {
	key := fmt.Sprintf(orderByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[orderMerchantCacheResponseDB](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderByMerchantCache(ctx context.Context, req *requests.FindAllOrderByMerchant, data []*db.GetOrdersByMerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetOrdersByMerchantRow{}
	}

	key := fmt.Sprintf(orderByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &orderMerchantCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}
