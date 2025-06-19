package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	categoryStatsMonthTotalPriceCacheKey = "category:stats:month:%d:year:%d"
	categoryStatsYearTotalPriceCacheKey  = "category:stats:year:%d"

	categoryStatsMonthPriceCacheKey = "category:stats:month:%d"
	categoryStatsYearPriceCacheKey  = "category:stats:year:%d"
)

type categoryStatsCache struct {
	store *CacheStore
}

func NewCategoryStatsCache(store *CacheStore) *categoryStatsCache {
	return &categoryStatsCache{store: store}
}

func (s *categoryStatsCache) GetCachedMonthTotalPriceCache(req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)

	result, found := GetFromCache[[]*response.CategoriesMonthlyTotalPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthTotalPriceCache(req *requests.MonthTotalPrice, data []*response.CategoriesMonthlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearTotalPriceCache(year int) ([]*response.CategoriesYearlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	result, found := GetFromCache[[]*response.CategoriesYearlyTotalPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedYearTotalPriceCache(year int, data []*response.CategoriesYearlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedMonthPriceCache(year int) ([]*response.CategoryMonthPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	result, found := GetFromCache[[]*response.CategoryMonthPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthPriceCache(year int, data []*response.CategoryMonthPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearPriceCache(year int) ([]*response.CategoryYearPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	result, found := GetFromCache[[]*response.CategoryYearPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsCache) SetCachedYearPriceCache(year int, data []*response.CategoryYearPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	SetToCache(s.store, key, &data, ttlDefault)
}
