package order_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

const (
	monthlyTotalRevenueCacheKeyByMerchant = "order:monthly:totalRevenue:merchant:%d:month:%d:year:%d"
	yearlyTotalRevenueCacheKeyByMerchant  = "order:yearly:totalRevenue:merchant:%d:year:%d"
	monthlyOrderCacheKeyByMerchant        = "order:monthly:order:merchant:%d:year:%d"
	yearlyOrderCacheKeyByMerchant         = "order:yearly:order:merchant:%d:year:%d"
)

type orderStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewOrderStatsByMerchantCache(store *cache.CacheStore) *orderStatsByMerchantCache {
	return &orderStatsByMerchantCache{store: store}
}

func (s *orderStatsByMerchantCache) GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) (*response.ApiResponseOrderMonthlyTotalRevenue, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseOrderMonthlyTotalRevenue](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, data *response.ApiResponseOrderMonthlyTotalRevenue) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) (*response.ApiResponseOrderYearlyTotalRevenue, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseOrderYearlyTotalRevenue](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderStatsByMerchantCache) SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, data *response.ApiResponseOrderYearlyTotalRevenue) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) (*response.ApiResponseOrderMonthly, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseOrderMonthly](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, data *response.ApiResponseOrderMonthly) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) (*response.ApiResponseOrderYearly, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[response.ApiResponseOrderYearly](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *orderStatsByMerchantCache) SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, data *response.ApiResponseOrderYearly) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
