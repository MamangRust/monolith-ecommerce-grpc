package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, bool)
	SetCachedCategoriesCache(req *requests.FindAllCategory, data []*response.CategoryResponse, total *int)

	GetCachedCategoryActiveCache(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, bool)
	SetCachedCategoryActiveCache(req *requests.FindAllCategory, data []*response.CategoryResponseDeleteAt, total *int)

	GetCachedCategoryTrashedCache(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, bool)
	SetCachedCategoryTrashedCache(req *requests.FindAllCategory, data []*response.CategoryResponseDeleteAt, total *int)

	GetCachedCategoryCache(id int) (*response.CategoryResponse, bool)
	SetCachedCategoryCache(data *response.CategoryResponse)
}

type CategoryCommandCache interface {
	DeleteCachedCategoryCache(id int)
}

type CategoryStatsCache interface {
	GetCachedMonthTotalPriceCache(req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, bool)
	SetCachedMonthTotalPriceCache(req *requests.MonthTotalPrice, data []*response.CategoriesMonthlyTotalPriceResponse)

	GetCachedYearTotalPriceCache(year int) ([]*response.CategoriesYearlyTotalPriceResponse, bool)
	SetCachedYearTotalPriceCache(year int, data []*response.CategoriesYearlyTotalPriceResponse)

	GetCachedMonthPriceCache(year int) ([]*response.CategoryMonthPriceResponse, bool)
	SetCachedMonthPriceCache(year int, data []*response.CategoryMonthPriceResponse)

	GetCachedYearPriceCache(year int) ([]*response.CategoryYearPriceResponse, bool)
	SetCachedYearPriceCache(year int, data []*response.CategoryYearPriceResponse)
}

type CategoryStatsByIdCache interface {
	GetCachedMonthTotalPriceByIdCache(req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, bool)
	SetCachedMonthTotalPriceByIdCache(req *requests.MonthTotalPriceCategory, data []*response.CategoriesMonthlyTotalPriceResponse)

	GetCachedYearTotalPriceByIdCache(req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, bool)
	SetCachedYearTotalPriceByIdCache(req *requests.YearTotalPriceCategory, data []*response.CategoriesYearlyTotalPriceResponse)

	GetCachedMonthPriceByIdCache(req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, bool)
	SetCachedMonthPriceByIdCache(req *requests.MonthPriceId, data []*response.CategoryMonthPriceResponse)

	GetCachedYearPriceByIdCache(req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, bool)
	SetCachedYearPriceByIdCache(req *requests.YearPriceId, data []*response.CategoryYearPriceResponse)
}

type CategoryStatsByMerchantCache interface {
	GetCachedMonthTotalPriceByMerchantCache(req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, bool)
	SetCachedMonthTotalPriceByMerchantCache(req *requests.MonthTotalPriceMerchant, data []*response.CategoriesMonthlyTotalPriceResponse)

	GetCachedYearTotalPriceByMerchantCache(req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, bool)
	SetCachedYearTotalPriceByMerchantCache(req *requests.YearTotalPriceMerchant, data []*response.CategoriesYearlyTotalPriceResponse)

	GetCachedMonthPriceByMerchantCache(req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, bool)
	SetCachedMonthPriceByMerchantCache(req *requests.MonthPriceMerchant, data []*response.CategoryMonthPriceResponse)

	GetCachedYearPriceByMerchantCache(req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, bool)
	SetCachedYearPriceByMerchantCache(req *requests.YearPriceMerchant, data []*response.CategoryYearPriceResponse)
}
