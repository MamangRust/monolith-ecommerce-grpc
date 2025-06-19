package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	orderAllCacheKey     = "order:all:page:%d:pageSize:%d:search:%s"
	orderByIdCacheKey    = "order:id:%d"
	orderActiveCacheKey  = "order:active:page:%d:pageSize:%d:search:%s"
	orderTrashedCacheKey = "order:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type orderCacheResponse struct {
	Data         []*response.OrderResponse `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type orderCacheResponseDeleteAt struct {
	Data         []*response.OrderResponseDeleteAt `json:"data"`
	TotalRecords *int                              `json:"total_records"`
}

type orderQueryCache struct {
	store *CacheStore
}

func NewOrderQueryCache(store *CacheStore) *orderQueryCache {
	return &orderQueryCache{store: store}
}

func (s *orderQueryCache) GetOrderAllCache(req *requests.FindAllOrder) ([]*response.OrderResponse, *int, bool) {
	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[orderCacheResponse](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderAllCache(req *requests.FindAllOrder, data []*response.OrderResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.OrderResponse{}
	}

	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponse{Data: data, TotalRecords: total}
	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetOrderActiveCache(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[orderCacheResponseDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderActiveCache(req *requests.FindAllOrder, data []*response.OrderResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.OrderResponseDeleteAt{}
	}

	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetOrderTrashedCache(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[orderCacheResponseDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderTrashedCache(req *requests.FindAllOrder, data []*response.OrderResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.OrderResponseDeleteAt{}
	}

	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderCache(order_id int) (*response.OrderResponse, bool) {
	key := fmt.Sprintf(orderByIdCacheKey, order_id)

	result, found := GetFromCache[*response.OrderResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderQueryCache) SetCachedOrderCache(data *response.OrderResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(orderByIdCacheKey, data.ID)
	SetToCache(s.store, key, data, ttlDefault)
}
