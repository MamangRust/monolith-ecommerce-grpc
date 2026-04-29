package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
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

func (s *orderStatsCache) GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, data []*db.GetMonthlyTotalRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTotalRevenueRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetYearlyTotalRevenueCache(ctx context.Context, year int, data []*db.GetYearlyTotalRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetMonthlyOrderCache(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyOrderRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetMonthlyOrderCache(ctx context.Context, year int, data []*db.GetMonthlyOrderRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyOrderCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyOrderCache(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKey, year)
	result, found := cache.GetFromCache[[]*db.GetYearlyOrderRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetYearlyOrderCache(ctx context.Context, year int, data []*db.GetYearlyOrderRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyOrderCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
