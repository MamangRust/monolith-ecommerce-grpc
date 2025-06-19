package mencache

import (
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

func (s *categoryStatsByMerchantCache) GetCachedMonthTotalPriceByMerchantCache(req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantMonthTotalPriceCacheKey, req.MerchantID, req.Month, req.Year)

	result, found := GetFromCache[[]*response.CategoriesMonthlyTotalPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedMonthTotalPriceByMerchantCache(req *requests.MonthTotalPriceMerchant, data []*response.CategoriesMonthlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantMonthTotalPriceCacheKey, req.MerchantID, req.Month, req.Year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedYearTotalPriceByMerchantCache(req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantYearTotalPriceCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.CategoriesYearlyTotalPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedYearTotalPriceByMerchantCache(req *requests.YearTotalPriceMerchant, data []*response.CategoriesYearlyTotalPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantYearTotalPriceCacheKey, req.MerchantID, req.Year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedMonthPriceByMerchantCache(req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantMonthPriceCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.CategoryMonthPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedMonthPriceByMerchantCache(req *requests.MonthPriceMerchant, data []*response.CategoryMonthPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantMonthPriceCacheKey, req.MerchantID, req.Year)

	SetToCache(s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByMerchantCache) GetCachedYearPriceByMerchantCache(req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, bool) {
	key := fmt.Sprintf(categoryStatsByMerchantYearPriceCacheKey, req.MerchantID, req.Year)

	result, found := GetFromCache[[]*response.CategoryYearPriceResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *categoryStatsByMerchantCache) SetCachedYearPriceByMerchantCache(req *requests.YearPriceMerchant, data []*response.CategoryYearPriceResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(categoryStatsByMerchantYearPriceCacheKey, req.MerchantID, req.Year)

	SetToCache(s.store, key, &data, ttlDefault)
}
