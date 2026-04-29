package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
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

func (s *categoryStatsCache) GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, bool) {
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalPriceRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data []*db.GetMonthlyTotalPriceRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, bool) {
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTotalPriceRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedYearTotalPriceCache(ctx context.Context, year int, data []*db.GetYearlyTotalPriceRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedMonthPriceCache(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, bool) {
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyCategoryRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthPriceCache(ctx context.Context, year int, data []*db.GetMonthlyCategoryRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearPriceCache(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, bool) {
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyCategoryRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedYearPriceCache(ctx context.Context, year int, data []*db.GetYearlyCategoryRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
