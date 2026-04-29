package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
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

func (s *orderStatsByMerchantCache) GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyTotalRevenueByMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, data []*db.GetMonthlyTotalRevenueByMerchantRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyTotalRevenueByMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, data []*db.GetYearlyTotalRevenueByMerchantRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyTotalRevenueCacheKeyByMerchant, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetMonthlyOrderByMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, data []*db.GetMonthlyOrderByMerchantRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsByMerchantCache) GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	result, found := cache.GetFromCache[[]*db.GetYearlyOrderByMerchantRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *orderStatsByMerchantCache) SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, data []*db.GetYearlyOrderByMerchantRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyOrderCacheKeyByMerchant, req.MerchantID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
