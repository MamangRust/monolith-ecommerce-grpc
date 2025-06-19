package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	orderItemAllCacheKey     = "order_item:all:page:%d:pageSize:%d:search:%s"
	orderItemActiveCacheKey  = "order_item:active:page:%d:pageSize:%d:search:%s"
	orderItemTrashedCacheKey = "order_item:trashed:page:%d:pageSize:%d:search:%s"

	orderItemByIdCacheKey = "order_item:id:%d"

	ttlDefault = 5 * time.Minute
)

type orderItemQueryCacheResponse struct {
	Data         []*response.OrderItemResponse `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type orderItemQueryCacheResponseDeleteAt struct {
	Data         []*response.OrderItemResponseDeleteAt `json:"data"`
	TotalRecords *int                                  `json:"total_records"`
}

type orderItemQueryCache struct {
	store *CacheStore
}

func NewOrderItemQueryCache(store *CacheStore) *orderItemQueryCache {
	return &orderItemQueryCache{store: store}
}

func (o *orderItemQueryCache) GetCachedOrderItemsAll(req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, bool) {
	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[orderItemQueryCacheResponse](o.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (O *orderItemQueryCache) SetCachedOrderItemsAll(req *requests.FindAllOrderItems, data []*response.OrderItemResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.OrderItemResponse{}
	}

	key := fmt.Sprintf(orderItemAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &orderItemQueryCacheResponse{Data: data, TotalRecords: total}
	SetToCache(O.store, key, payload, ttlDefault)
}

func (O *orderItemQueryCache) GetCachedOrderItemActive(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[orderItemQueryCacheResponseDeleteAt](O.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (O *orderItemQueryCache) SetCachedOrderItemActive(req *requests.FindAllOrderItems, data []*response.OrderItemResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.OrderItemResponseDeleteAt{}
	}

	key := fmt.Sprintf(orderItemActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderItemQueryCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(O.store, key, payload, ttlDefault)
}

func (O *orderItemQueryCache) GetCachedOrderItemTrashed(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[orderItemQueryCacheResponseDeleteAt](O.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (O *orderItemQueryCache) SetCachedOrderItemTrashed(req *requests.FindAllOrderItems, data []*response.OrderItemResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.OrderItemResponseDeleteAt{}
	}

	key := fmt.Sprintf(orderItemTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderItemQueryCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(O.store, key, payload, ttlDefault)
}

func (O *orderItemQueryCache) GetCachedOrderItems(order_id int) ([]*response.OrderItemResponse, bool) {
	key := fmt.Sprintf(orderItemByIdCacheKey, order_id)
	result, found := GetFromCache[[]*response.OrderItemResponse](O.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (O *orderItemQueryCache) SetCachedOrderItems(data []*response.OrderItemResponse) {
	if len(data) == 0 {
		return
	}

	key := fmt.Sprintf(orderItemByIdCacheKey, data[0].OrderID)
	SetToCache(O.store, key, &data, ttlDefault)
}
