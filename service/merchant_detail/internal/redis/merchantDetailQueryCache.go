package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	merchantDetailAllCacheKey     = "merchant_detail:all:page:%d:pageSize:%d:search:%s"
	merchantDetailByIdCacheKey    = "merchant_detail:id:%d"
	merchantDetailActiveCacheKey  = "merchant_detail:active:page:%d:pageSize:%d:search:%s"
	merchantDetailTrashedCacheKey = "merchant_detail:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantDetailCacheResponse struct {
	Data         []*response.MerchantDetailResponse `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type merchantDetailCacheResponseDeleteAt struct {
	Data         []*response.MerchantDetailResponseDeleteAt `json:"data"`
	TotalRecords *int                                       `json:"total_records"`
}

type merchantDetailQueryCache struct {
	store *CacheStore
}

func NewMerchantDetailQueryCache(store *CacheStore) *merchantDetailQueryCache {
	return &merchantDetailQueryCache{
		store: store,
	}
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailAll(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponse, *int, bool) {
	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDetailCacheResponse](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailAll(req *requests.FindAllMerchant, data []*response.MerchantDetailResponse, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantDetailResponse{}
	}

	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailCacheResponse{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailActive(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDetailCacheResponseDeleteAt](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailActive(req *requests.FindAllMerchant, data []*response.MerchantDetailResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantDetailResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailCacheResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailTrashed(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDetailCacheResponseDeleteAt](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailTrashed(req *requests.FindAllMerchant, data []*response.MerchantDetailResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantDetailResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailCacheResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetail(id int) (*response.MerchantDetailResponse, bool) {
	key := fmt.Sprintf(merchantDetailByIdCacheKey, id)

	result, found := GetFromCache[*response.MerchantDetailResponse](m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetail(data *response.MerchantDetailResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDetailByIdCacheKey, data.ID)
	SetToCache(m.store, key, data, ttlDefault)
}
