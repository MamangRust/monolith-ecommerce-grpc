package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	categoryAllCacheKey     = "category:all:page:%d:pageSize:%d:search:%s"
	categoryByIdCacheKey    = "category:id:%d"
	categoryActiveCacheKey  = "category:active:page:%d:pageSize:%d:search:%s"
	categoryTrashedCacheKey = "category:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type categoryCacheResponse struct {
	Data         []*response.CategoryResponse `json:"data"`
	TotalRecords *int                         `json:"totalRecords"`
}

type categoryCacheResponseDeleteAt struct {
	Data         []*response.CategoryResponseDeleteAt `json:"data"`
	TotalRecords *int                                 `json:"totalRecords"`
}

type categoryQueryCache struct {
	store *CacheStore
}

func NewCategoryQueryCache(store *CacheStore) *categoryQueryCache {
	return &categoryQueryCache{store: store}
}

func (s *categoryQueryCache) GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, bool) {
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[categoryCacheResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*response.CategoryResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.CategoryResponse{}
	}

	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[categoryCacheResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}
func (s *categoryQueryCache) SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*response.CategoryResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CategoryResponseDeleteAt{}
	}

	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[categoryCacheResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*response.CategoryResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CategoryResponseDeleteAt{}
	}

	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryCache(ctx context.Context, id int) (*response.CategoryResponse, bool) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	result, found := GetFromCache[*response.CategoryResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryQueryCache) SetCachedCategoryCache(ctx context.Context, data *response.CategoryResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryByIdCacheKey, data.ID)

	SetToCache(ctx, s.store, key, data, ttlDefault)
}
