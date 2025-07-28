package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, bool)
	SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*response.CategoryResponse, total *int)

	GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, bool)
	SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*response.CategoryResponseDeleteAt, total *int)

	GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, bool)
	SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*response.CategoryResponseDeleteAt, total *int)

	GetCachedCategoryCache(ctx context.Context, id int) (*response.CategoryResponse, bool)
	SetCachedCategoryCache(ctx context.Context, data *response.CategoryResponse)
}
type CategoryCommandCache interface {
	DeleteCachedCategoryCache(ctx context.Context, id int)
}

type CategoryStatsCache interface {
	GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, bool)
	SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data []*response.CategoriesMonthlyTotalPriceResponse)

	GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*response.CategoriesYearlyTotalPriceResponse, bool)
	SetCachedYearTotalPriceCache(ctx context.Context, year int, data []*response.CategoriesYearlyTotalPriceResponse)

	GetCachedMonthPriceCache(ctx context.Context, year int) ([]*response.CategoryMonthPriceResponse, bool)
	SetCachedMonthPriceCache(ctx context.Context, year int, data []*response.CategoryMonthPriceResponse)

	GetCachedYearPriceCache(ctx context.Context, year int) ([]*response.CategoryYearPriceResponse, bool)
	SetCachedYearPriceCache(ctx context.Context, year int, data []*response.CategoryYearPriceResponse)
}

type CategoryStatsByIdCache interface {
	GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, bool)
	SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data []*response.CategoriesMonthlyTotalPriceResponse)

	GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, bool)
	SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data []*response.CategoriesYearlyTotalPriceResponse)

	GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, bool)
	SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data []*response.CategoryMonthPriceResponse)

	GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, bool)
	SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data []*response.CategoryYearPriceResponse)
}

type CategoryStatsByMerchantCache interface {
	GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, bool)
	SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant, data []*response.CategoriesMonthlyTotalPriceResponse)

	GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, bool)
	SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant, data []*response.CategoriesYearlyTotalPriceResponse)

	GetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, bool)
	SetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant, data []*response.CategoryMonthPriceResponse)

	GetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, bool)
	SetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant, data []*response.CategoryYearPriceResponse)
}
