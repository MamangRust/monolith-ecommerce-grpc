package order_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

const (
	monthlyTotalRevenueCacheKey = "order:monthly:totalRevenue:month:%d:year:%d"
	yearlyTotalRevenueCacheKey  = "order:yearly:totalRevenue:year:%d"
	monthlyOrderCacheKey        = "order:monthly:order:month:%d"
	yearlyOrderCacheKey         = "order:yearly:order:year:%d"
)

type orderStatsCache struct {
	store *cache.CacheStore
}

func NewOrderStatsCache(store *cache.CacheStore) *orderStatsCache {
	return &orderStatsCache{store: store}
}

func (s *orderStatsCache) GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) (*response.ApiResponseOrderMonthlyTotalRevenue, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseOrderMonthlyTotalRevenue](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderStatsCache) SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, data *response.ApiResponseOrderMonthlyTotalRevenue) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyTotalRevenueCache(ctx context.Context, year int) (*response.ApiResponseOrderYearlyTotalRevenue, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)
	result, found := cache.GetFromCache[response.ApiResponseOrderYearlyTotalRevenue](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderStatsCache) SetYearlyTotalRevenueCache(ctx context.Context, year int, data *response.ApiResponseOrderYearlyTotalRevenue) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderStatsCache) GetMonthlyOrderCache(ctx context.Context, year int) (*response.ApiResponseOrderMonthly, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKey, year)
	result, found := cache.GetFromCache[response.ApiResponseOrderMonthly](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderStatsCache) SetMonthlyOrderCache(ctx context.Context, year int, data *response.ApiResponseOrderMonthly) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyOrderCacheKey, year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyOrderCache(ctx context.Context, year int) (*response.ApiResponseOrderYearly, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKey, year)
	result, found := cache.GetFromCache[response.ApiResponseOrderYearly](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderStatsCache) SetYearlyOrderCache(ctx context.Context, year int, data *response.ApiResponseOrderYearly) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyOrderCacheKey, year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
