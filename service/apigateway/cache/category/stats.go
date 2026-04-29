package category_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

const (
	categoryStatsMonthTotalPriceCacheKey = "category:stats:month:%d:year:%d"
	categoryStatsYearTotalPriceCacheKey  = "category:stats:year:%d"
	categoryStatsMonthPriceCacheKey      = "category:stats:month:%d"
	categoryStatsYearPriceCacheKey       = "category:stats:year:%d"
)

type categoryStatsCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsCache(store *cache.CacheStore) *categoryStatsCache {
	return &categoryStatsCache{store: store}
}

func (s *categoryStatsCache) GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) (*response.ApiResponseCategoryMonthlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseCategoryMonthlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsCache) SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data *response.ApiResponseCategoryMonthlyTotalPrice) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearTotalPriceCache(ctx context.Context, year int) (*response.ApiResponseCategoryYearlyTotalPrice, bool) {
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	result, found := cache.GetFromCache[response.ApiResponseCategoryYearlyTotalPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsCache) SetCachedYearTotalPriceCache(ctx context.Context, year int, data *response.ApiResponseCategoryYearlyTotalPrice) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedMonthPriceCache(ctx context.Context, year int) (*response.ApiResponseCategoryMonthPrice, bool) {
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	result, found := cache.GetFromCache[response.ApiResponseCategoryMonthPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsCache) SetCachedMonthPriceCache(ctx context.Context, year int, data *response.ApiResponseCategoryMonthPrice) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearPriceCache(ctx context.Context, year int) (*response.ApiResponseCategoryYearPrice, bool) {
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	result, found := cache.GetFromCache[response.ApiResponseCategoryYearPrice](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *categoryStatsCache) SetCachedYearPriceCache(ctx context.Context, year int, data *response.ApiResponseCategoryYearPrice) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
