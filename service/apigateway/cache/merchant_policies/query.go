package merchantpolicies_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	merchantPolicyAllCacheKey     = "merchant_policy:all:page:%d:pageSize:%d:search:%s"
	merchantPolicyByIdCacheKey    = "merchant_policy:id:%d"
	merchantPolicyActiveCacheKey  = "merchant_policy:active:page:%d:pageSize:%d:search:%s"
	merchantPolicyTrashedCacheKey = "merchant_policy:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantPolicyQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantPolicyQueryCache(store *cache.CacheStore) *merchantPolicyQueryCache {
	return &merchantPolicyQueryCache{
		store: store,
	}
}

func (m *merchantPolicyQueryCache) GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantPolicies, bool) {
	key := fmt.Sprintf(merchantPolicyAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantPolicies](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantPolicyQueryCache) SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantPolicies) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantPolicyAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantPolicyQueryCache) GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantPoliciesDeleteAt, bool) {
	key := fmt.Sprintf(merchantPolicyActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantPoliciesDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantPolicyQueryCache) SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantPoliciesDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantPolicyActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantPolicyQueryCache) GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantPoliciesDeleteAt, bool) {
	key := fmt.Sprintf(merchantPolicyTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantPoliciesDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantPolicyQueryCache) SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantPoliciesDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantPolicyTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantPolicyQueryCache) GetCachedMerchantPolicy(ctx context.Context, id int) (*response.ApiResponseMerchantPolicies, bool) {
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, id)

	result, found := cache.GetFromCache[response.ApiResponseMerchantPolicies](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantPolicyQueryCache) SetCachedMerchantPolicy(ctx context.Context, data *response.ApiResponseMerchantPolicies) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantPolicyByIdCacheKey, data.Data.ID)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
