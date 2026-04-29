package category_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	categoryAllCacheKey     = "category:all:page:%d:pageSize:%d:search:%s"
	categoryByIdCacheKey    = "category:id:%d"
	categoryActiveCacheKey  = "category:active:page:%d:pageSize:%d:search:%s"
	categoryTrashedCacheKey = "category:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type categoryQueryCache struct {
	store *cache.CacheStore
}

func NewCategoryQueryCache(store *cache.CacheStore) *categoryQueryCache {
	return &categoryQueryCache{store: store}
}

func (s *categoryQueryCache) GetCachedCategoriesCache(
	ctx context.Context,
	req *requests.FindAllCategory,
) (*response.ApiResponsePaginationCategory, bool) {
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationCategory](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryQueryCache) SetCachedCategoriesCache(
	ctx context.Context,
	req *requests.FindAllCategory,
	data *response.ApiResponsePaginationCategory,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryActiveCache(
	ctx context.Context,
	req *requests.FindAllCategory,
) (*response.ApiResponsePaginationCategoryDeleteAt, bool) {
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationCategoryDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryQueryCache) SetCachedCategoryActiveCache(
	ctx context.Context,
	req *requests.FindAllCategory,
	data *response.ApiResponsePaginationCategoryDeleteAt,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryTrashedCache(
	ctx context.Context,
	req *requests.FindAllCategory,
) (*response.ApiResponsePaginationCategoryDeleteAt, bool) {
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationCategoryDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryQueryCache) SetCachedCategoryTrashedCache(
	ctx context.Context,
	req *requests.FindAllCategory,
	data *response.ApiResponsePaginationCategoryDeleteAt,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryCache(
	ctx context.Context,
	id int,
) (*response.ApiResponseCategory, bool) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	result, found := cache.GetFromCache[response.ApiResponseCategory](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryQueryCache) SetCachedCategoryCache(
	ctx context.Context,
	data *response.ApiResponseCategory,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryByIdCacheKey, data.Data.ID)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
