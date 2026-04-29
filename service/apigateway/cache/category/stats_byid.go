package category_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

const (
	categoryStatsByIdMonthTotalPriceCacheKey = "category:stats:byid:%d:month:%d:year:%d"
	categoryStatsByIdYearTotalPriceCacheKey  = "category:stats:byid:%d:year:%d"
	categoryStatsByIdMonthPriceCacheKey      = "category:stats:byid:%d:month:%d"
	categoryStatsByIdYearPriceCacheKey       = "category:stats:byid:%d:year:%d"
)

type categoryStatsByIdCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsByIdCache(store *cache.CacheStore) *categoryStatsByIdCache {
	return &categoryStatsByIdCache{store: store}
}

func (s *categoryStatsByIdCache) GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) (*response.ApiResponseCategoryMonthlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseCategoryMonthlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data *response.ApiResponseCategoryMonthlyTotalPrice) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) (*response.ApiResponseCategoryYearlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseCategoryYearlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByIdCache) SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data *response.ApiResponseCategoryYearlyTotalPrice) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) (*response.ApiResponseCategoryMonthPrice, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseCategoryMonthPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data *response.ApiResponseCategoryMonthPrice) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) (*response.ApiResponseCategoryYearPrice, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseCategoryYearPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsByIdCache) SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data *response.ApiResponseCategoryYearPrice) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
