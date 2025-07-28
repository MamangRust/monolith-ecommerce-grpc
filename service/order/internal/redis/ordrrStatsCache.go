package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	monthlyTotalRevenueCacheKey = "order:monthly:totalRevenue:month:%d:year:%d"
	yearlyTotalRevenueCacheKey  = "order:yearly:totalRevenue:year:%d"

	monthlyOrderCacheKey = "order:monthly:order:month:%d"
	yearlyOrderCacheKey  = "order:yearly:order:year:%d"
)

type orderStatsCache struct {
	store *CacheStore
}

func NewOrderStatsCache(store *CacheStore) *orderStatsCache {
	return &orderStatsCache{store: store}
}

func (s *orderStatsCache) GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.OrderMonthlyTotalRevenueResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, data []*response.OrderMonthlyTotalRevenueResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*response.OrderYearlyTotalRevenueResponse, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)

	result, found := GetFromCache[[]*response.OrderYearlyTotalRevenueResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetYearlyTotalRevenueCache(ctx context.Context, year int, data []*response.OrderYearlyTotalRevenueResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetMonthlyOrderCache(ctx context.Context, year int) ([]*response.OrderMonthlyResponse, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKey, year)

	result, found := GetFromCache[[]*response.OrderMonthlyResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetMonthlyOrderCache(ctx context.Context, year int, data []*response.OrderMonthlyResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthlyOrderCacheKey, year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyOrderCache(ctx context.Context, year int) ([]*response.OrderYearlyResponse, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKey, year)

	result, found := GetFromCache[[]*response.OrderYearlyResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetYearlyOrderCache(ctx context.Context, year int, data []*response.OrderYearlyResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearlyOrderCacheKey, year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}
