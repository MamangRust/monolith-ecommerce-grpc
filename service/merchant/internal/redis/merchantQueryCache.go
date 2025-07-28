package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	merchantAllCacheKey      = "merchant:all:page:%d:pageSize:%d:search:%s"
	merchantByIdCacheKey     = "merchant:id:%d"
	merchantActiveCacheKey   = "merchant:active:page:%d:pageSize:%d:search:%s"
	merchantTrashedCacheKey  = "merchant:trashed:page:%d:pageSize:%d:search:%s"
	merchantByApiKeyCacheKey = "merchant:api_key:%s"
	merchantByUserIdCacheKey = "merchant:user_id:%d"

	ttlDefault = 5 * time.Minute
)

type merchantCachedResponse struct {
	Data         []*response.MerchantResponse `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type merchantCachedResponseDeleteAt struct {
	Data         []*response.MerchantResponseDeleteAt `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantQueryCache struct {
	store *CacheStore
}

func NewMerchantQueryCache(store *CacheStore) *merchantQueryCache {
	return &merchantQueryCache{store: store}
}

func (m *merchantQueryCache) GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponse, *int, bool) {
	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantCachedResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantResponse{}
	}

	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchant(ctx context.Context, id int) (*response.MerchantResponse, bool) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	result, found := GetFromCache[*response.MerchantResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantQueryCache) SetCachedMerchant(ctx context.Context, data *response.MerchantResponse) {
	key := fmt.Sprintf(merchantByIdCacheKey, data.ID)

	SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantsByUserId(ctx context.Context, id int) ([]*response.MerchantResponse, bool) {
	key := fmt.Sprintf(merchantByUserIdCacheKey, id)

	result, found := GetFromCache[[]*response.MerchantResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantQueryCache) SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*response.MerchantResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByUserIdCacheKey, userId)

	SetToCache(ctx, m.store, key, &data, ttlDefault)
}
