package mencache

import (
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

func (s *orderStatsCache) GetMonthlyTotalRevenueCache(req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.OrderMonthlyTotalRevenueResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetMonthlyTotalRevenueCache(req *requests.MonthTotalRevenue, data []*response.OrderMonthlyTotalRevenueResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyTotalRevenueCache(year int) ([]*response.OrderYearlyTotalRevenueResponse, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)

	result, found := GetFromCache[[]*response.OrderYearlyTotalRevenueResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetYearlyTotalRevenueCache(year int, data []*response.OrderYearlyTotalRevenueResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetMonthlyOrderCache(year int) ([]*response.OrderMonthlyResponse, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKey, year)

	result, found := GetFromCache[[]*response.OrderMonthlyResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetMonthlyOrderCache(year int, data []*response.OrderMonthlyResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthlyOrderCacheKey, year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyOrderCache(year int) ([]*response.OrderYearlyResponse, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKey, year)

	result, found := GetFromCache[[]*response.OrderYearlyResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsCache) SetYearlyOrderCache(year int, data []*response.OrderYearlyResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearlyOrderCacheKey, year)
	SetToCache(s.store, key, &data, ttlDefault)
}
