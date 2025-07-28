package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	categoryStatsByMerchantMonthTotalPriceCacheKey = "category:stats:bymerchant:%d:month:%d:year:%d"
	categoryStatsByMerchantYearTotalPriceCacheKey  = "category:stats:bymerchant:%d:year:%d"

	categoryStatsByMerchantMonthPriceCacheKey = "category:stats:bymerchant:%d:month:%d"
	categoryStatsByMerchantYearPriceCacheKey  = "category:stats:bymerchant:%d:year:%d"
)

type categoryStatsByMerchantCache struct {
	store *CacheStore
}

func NewCategoryStatsByMerchantCache(store *CacheStore) *categoryStatsByMerchantCache {
	return &categoryStatsByMerchantCache{store: store}
}

func (s *categoryStatsByMerchantCache) GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantMonthTotalPriceCacheKey, req.MerchantID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.CategoriesMonthlyTotalPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant, data []*response.CategoriesMonthlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantMonthTotalPriceCacheKey, req.MerchantID, req.Month, req.Year)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantYearTotalPriceCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.CategoriesYearlyTotalPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant, data []*response.CategoriesYearlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantYearTotalPriceCacheKey, req.MerchantID, req.Year)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantMonthPriceCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.CategoryMonthPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant, data []*response.CategoryMonthPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantMonthPriceCacheKey, req.MerchantID, req.Year)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantYearPriceCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.CategoryYearPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant, data []*response.CategoryYearPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantYearPriceCacheKey, req.MerchantID, req.Year)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}
