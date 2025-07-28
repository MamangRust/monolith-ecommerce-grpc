package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	categoryStatsByIdMonthTotalPriceCacheKey = "category:stats:byid:%d:month:%d:year:%d"
	categoryStatsByIdYearTotalPriceCacheKey  = "category:stats:byid:%d:year:%d"

	categoryStatsByIdMonthPriceCacheKey = "category:stats:byid:%d:month:%d"
	categoryStatsByIdYearPriceCacheKey  = "category:stats:byid:%d:year:%d"
)

type categoryStatsByIdCache struct {
	store *CacheStore
}

func NewCategoryStatsByIdCache(store *CacheStore) *categoryStatsByIdCache {
	return &categoryStatsByIdCache{store: store}
}

func (s *categoryStatsByIdCache) GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.CategoriesMonthlyTotalPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data []*response.CategoriesMonthlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)

	result, found := GetFromCache[[]*response.CategoriesYearlyTotalPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data []*response.CategoriesYearlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)

	result, found := GetFromCache[[]*response.CategoryMonthPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data []*response.CategoryMonthPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)

	SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)

	result, found := GetFromCache[[]*response.CategoryYearPriceResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data []*response.CategoryYearPriceResponse) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)
	SetToCache(ctx, s.store, key, &data, ttlDefault)
}
