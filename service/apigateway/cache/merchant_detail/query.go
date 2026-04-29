package merchantdetail_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	merchantDetailAllCacheKey          = "merchant_detail:all:page:%d:pageSize:%d:search:%s"
	merchantDetailByIdCacheKey         = "merchant_detail:id:%d"
	merchantDetailActiveCacheKey       = "merchant_detail:active:page:%d:pageSize:%d:search:%s"
	merchantDetailTrashedCacheKey      = "merchant_detail:trashed:page:%d:pageSize:%d:search:%s"
	merchantDetailRelationByIdCacheKey = "merchant:detail:relation:id:%d"

	ttlDefault = 5 * time.Minute
)

type merchantDetailQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantDetailQueryCache(store *cache.CacheStore) *merchantDetailQueryCache {
	return &merchantDetailQueryCache{
		store: store,
	}
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetail, bool) {
	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantDetail](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetail) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetailDeleteAt, bool) {
	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantDetailDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetailDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetailDeleteAt, bool) {
	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationMerchantDetailDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetailDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetail(ctx context.Context, id int) (*response.ApiResponseMerchantDetail, bool) {
	key := fmt.Sprintf(merchantDetailByIdCacheKey, id)

	result, found := cache.GetFromCache[response.ApiResponseMerchantDetail](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetail(ctx context.Context, data *response.ApiResponseMerchantDetail) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDetailByIdCacheKey, data.Data.ID)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailRelation(
	ctx context.Context,
	merchantID int,
) (*response.ApiResponseMerchantDetailRelation, bool) {

	key := fmt.Sprintf(merchantDetailRelationByIdCacheKey, merchantID)

	result, found := cache.GetFromCache[response.ApiResponseMerchantDetailRelation](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailRelation(
	ctx context.Context,
	merchantID int,
	data *response.ApiResponseMerchantDetailRelation,
) {
	if merchantID <= 0 || data == nil || data.Data.ID != merchantID {
		return
	}

	key := fmt.Sprintf(merchantDetailRelationByIdCacheKey, merchantID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
