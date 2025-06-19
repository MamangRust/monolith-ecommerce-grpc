package mencache

import (
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

func (s *categoryStatsByIdCache) GetCachedMonthTotalPriceByIdCache(req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.CategoriesMonthlyTotalPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthTotalPriceByIdCache(req *requests.MonthTotalPriceCategory, data []*response.CategoriesMonthlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearTotalPriceByIdCache(req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)

	result, found := GetFromCache[[]*response.CategoriesYearlyTotalPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedYearTotalPriceByIdCache(req *requests.YearTotalPriceCategory, data []*response.CategoriesYearlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedMonthPriceByIdCache(req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)

	result, found := GetFromCache[[]*response.CategoryMonthPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthPriceByIdCache(req *requests.MonthPriceId, data []*response.CategoryMonthPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearPriceByIdCache(req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)

	result, found := GetFromCache[[]*response.CategoryYearPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedYearPriceByIdCache(req *requests.YearPriceId, data []*response.CategoryYearPriceResponse) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)
	SetToCache(s.store, key, &data, ttlDefault)
}
