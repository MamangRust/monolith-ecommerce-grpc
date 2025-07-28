package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	merchantBusinessAllCacheKey     = "merchant_business:all:page:%d:pageSize:%d:search:%s"
	merchantBusinessByIdCacheKey    = "merchant_business:id:%d"
	merchantBusinessActiveCacheKey  = "merchant_business:active:page:%d:pageSize:%d:search:%s"
	merchantBusinessTrashedCacheKey = "merchant_business:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantBusinessCacheResponse struct {
	Data         []*response.MerchantBusinessResponse `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantBusinessCacheResponseDeleteAt struct {
	Data         []*response.MerchantBusinessResponseDeleteAt `json:"data"`
	TotalRecords *int                                         `json:"total_records"`
}

type merchantBusinessQueryCache struct {
	store *CacheStore
}

func NewMerchantBusinessQueryCache(store *CacheStore) *merchantBusinessQueryCache {
	return &merchantBusinessQueryCache{
		store: store,
	}
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponse, *int, bool) {
	key := fmt.Sprintf(merchantBusinessAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantBusinessCacheResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantBusinessResponse, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantBusinessResponse{}
	}

	key := fmt.Sprintf(merchantBusinessAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessCacheResponse{Data: data, TotalRecords: totalRecords}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantBusinessActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantBusinessCacheResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantBusinessResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantBusinessResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantBusinessActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessCacheResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantBusinessTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantBusinessCacheResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantBusinessResponseDeleteAt, totalRecords *int) {
	if totalRecords == nil {
		zero := 0
		totalRecords = &zero
	}

	if data == nil {
		data = []*response.MerchantBusinessResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantBusinessTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessCacheResponseDeleteAt{Data: data, TotalRecords: totalRecords}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusiness(ctx context.Context, id int) (*response.MerchantBusinessResponse, bool) {
	key := fmt.Sprintf(merchantBusinessByIdCacheKey, id)

	result, found := GetFromCache[*response.MerchantBusinessResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusiness(ctx context.Context, data *response.MerchantBusinessResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantBusinessByIdCacheKey, data.ID)
	SetToCache(ctx, m.store, key, data, ttlDefault)
}
