package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	monthlyTotalRevenueCacheKeyByMerchant = "order:monthly:totalRevenue:merchant:%d:month:%d:year:%d"
	yearlyTotalRevenueCacheKeyByMerchant  = "order:yearly:totalRevenue:merchant:%d:year:%d"

	monthlyOrderCacheKeyByMerchant = "order:monthly:order:merchant:%d:year:%d"
	yearlyOrderCacheKeyByMerchant  = "order:yearly:order:merchant:%d:year:%d"
)

type orderStatsByMerchantCache struct {
	store *CacheStore
}

func NewOrderStatsByMerchantCache(store *CacheStore) *orderStatsByMerchantCache {
	return &orderStatsByMerchantCache{store: store}
}

func (s *orderStatsByMerchantCache) GetMonthlyTotalRevenueByMerchantCache(req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.OrderMonthlyTotalRevenueResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyTotalRevenueByMerchantCache(req *requests.MonthTotalRevenueMerchant, data []*response.OrderMonthlyTotalRevenueResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyTotalRevenueByMerchantCache(req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.OrderYearlyTotalRevenueResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetYearlyTotalRevenueByMerchantCache(req *requests.YearTotalRevenueMerchant, data []*response.OrderYearlyTotalRevenueResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetMonthlyOrderByMerchantCache(req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.OrderMonthlyResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyOrderByMerchantCache(req *requests.MonthOrderMerchant, data []*response.OrderMonthlyResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyOrderByMerchantCache(req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.OrderYearlyResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetYearlyOrderByMerchantCache(req *requests.YearOrderMerchant, data []*response.OrderYearlyResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}
