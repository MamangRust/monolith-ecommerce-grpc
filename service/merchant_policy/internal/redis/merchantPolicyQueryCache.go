package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	merchantPolicyAllCacheKey  = "merchant_policy:all:page:%d:pageSize:%d:search:%s"
	merchantPolicyByIdCacheKey = "merchant_policy:id:%d"

	merchantPolicyActiveCacheKey  = "merchant_policy:active:page:%d:pageSize:%d:search:%s"
	merchantPolicyTrashedCacheKey = "merchant_policy:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantPolicyCacheResponse struct {
	Data         []*response.MerchantPoliciesResponse `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantPolicyCacheResponseDeleteAt struct {
	Data         []*response.MerchantPoliciesResponseDeleteAt `json:"data"`
	TotalRecords *int                                         `json:"total_records"`
}

type merchantPolicyQueryCache struct {
	store *CacheStore
}

func NewMerchantPolicyQueryCache(store *CacheStore) *merchantPolicyQueryCache {
	return &merchantPolicyQueryCache{
		store: store,
	}
}

func (m *merchantPolicyQueryCache) GetCachedMerchantPolicyAll(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, bool) {
	key := fmt.Sprintf(merchantPolicyAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantPolicyCacheResponse](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantPolicyQueryCache) SetCachedMerchantPolicyAll(req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponse, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantPoliciesResponse{}
	}

	key := fmt.Sprintf(merchantPolicyAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantPolicyCacheResponse{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantPolicyQueryCache) GetCachedMerchantPolicyActive(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantPolicyActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantPolicyCacheResponseDeleteAt](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantPolicyQueryCache) SetCachedMerchantPolicyActive(req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantPoliciesResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantPolicyActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantPolicyCacheResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantPolicyQueryCache) GetCachedMerchantPolicyTrashed(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantPolicyTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantPolicyCacheResponseDeleteAt](m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantPolicyQueryCache) SetCachedMerchantPolicyTrashed(req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantPoliciesResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantPolicyTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantPolicyCacheResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(m.store, key, payload, ttlDefault)
}

func (m *merchantPolicyQueryCache) GetCachedMerchantPolicy(id int) (*response.MerchantPoliciesResponse, bool) {
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, id)

	result, found := GetFromCache[*response.MerchantPoliciesResponse](m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantPolicyQueryCache) SetCachedMerchantPolicy(data *response.MerchantPoliciesResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantPolicyByIdCacheKey, data.ID)
	SetToCache(m.store, key, data, ttlDefault)
}
