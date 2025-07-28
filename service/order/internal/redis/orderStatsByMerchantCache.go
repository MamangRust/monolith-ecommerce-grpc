package mencache

import (
	"context"
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

func (s *orderStatsByMerchantCache) GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.OrderMonthlyTotalRevenueResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, data []*response.OrderMonthlyTotalRevenueResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.OrderYearlyTotalRevenueResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, data []*response.OrderYearlyTotalRevenueResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.OrderMonthlyResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, data []*response.OrderMonthlyResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.OrderYearlyResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, data []*response.OrderYearlyResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}
