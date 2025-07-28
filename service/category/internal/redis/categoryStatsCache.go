package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	categoryStatsMonthTotalPriceCacheKey = "category:stats:month:%d:year:%d"
	categoryStatsYearTotalPriceCacheKey  = "category:stats:year:%d"

	categoryStatsMonthPriceCacheKey = "category:stats:month:%d"
	categoryStatsYearPriceCacheKey  = "category:stats:year:%d"
)

type categoryStatsCache struct {
	store *CacheStore
}

func NewCategoryStatsCache(store *CacheStore) *categoryStatsCache {
	return &categoryStatsCache{store: store}
}

func (s *categoryStatsCache) GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.CategoriesMonthlyTotalPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data []*response.CategoriesMonthlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*response.CategoriesYearlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	result, found := GetFromCache[[]*response.CategoriesYearlyTotalPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedYearTotalPriceCache(ctx context.Context, year int, data []*response.CategoriesYearlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedMonthPriceCache(ctx context.Context, year int) ([]*response.CategoryMonthPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	result, found := GetFromCache[[]*response.CategoryMonthPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthPriceCache(ctx context.Context, year int, data []*response.CategoryMonthPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearPriceCache(ctx context.Context, year int) ([]*response.CategoryYearPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	result, found := GetFromCache[[]*response.CategoryYearPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedYearPriceCache(ctx context.Context, year int, data []*response.CategoryYearPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}
