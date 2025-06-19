package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	merchantAwardAllCacheKey     = "merchant_award:all:page:%d:pageSize:%d:search:%s"
	merchantAwardByIdCacheKey    = "merchant_award:id:%d"
	merchantAwardActiveCacheKey  = "merchant_award:active:page:%d:pageSize:%d:search:%s"
	merchantAwardTrashedCacheKey = "merchant_award:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantAwardCacheResponse struct {
	Data         []*response.MerchantAwardResponse `json:"data"`
	TotalRecords *int                              `json:"total_records"`
}

type merchantAwardCacheResponseDeleteAt struct {
	Data         []*response.MerchantAwardResponseDeleteAt `json:"data"`
	TotalRecords *int                                      `json:"total_records"`
}

type merchantAwardQueryCache struct {
	store *CacheStore
}

func NewMerchantAwardQueryCache(store *CacheStore) *merchantAwardQueryCache {
	return &merchantAwardQueryCache{
		store: store,
	}
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardAll(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, bool) {
	key := fmt.Sprintf(merchantAwardAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantAwardCacheResponse](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardAll(req *requests.FindAllMerchant, data []*response.MerchantAwardResponse, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantAwardResponse{}
	}

	key := fmt.Sprintf(merchantAwardAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantAwardCacheResponse{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardActive(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantAwardActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantAwardCacheResponseDeleteAt](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardActive(req *requests.FindAllMerchant, data []*response.MerchantAwardResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}
	key := fmt.Sprintf(merchantAwardActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantAwardCacheResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardTrashed(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantAwardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantAwardCacheResponseDeleteAt](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardTrashed(req *requests.FindAllMerchant, data []*response.MerchantAwardResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantAwardResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantAwardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantAwardCacheResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAward(id int) (*response.MerchantAwardResponse, bool) {
	key := fmt.Sprintf(merchantAwardByIdCacheKey, id)

	result, found := GetFromCache[*response.MerchantAwardResponse](m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAward(data *response.MerchantAwardResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantAwardByIdCacheKey, data.ID)
	SetToCache(m.store, key, data, ttlDefault)
}
