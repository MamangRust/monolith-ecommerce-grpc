package category_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) (*response.ApiResponsePaginationCategory, bool)
	SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data *response.ApiResponsePaginationCategory)
	GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) (*response.ApiResponsePaginationCategoryDeleteAt, bool)
	SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data *response.ApiResponsePaginationCategoryDeleteAt)
	GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) (*response.ApiResponsePaginationCategoryDeleteAt, bool)
	SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data *response.ApiResponsePaginationCategoryDeleteAt)
	GetCachedCategoryCache(ctx context.Context, id int) (*response.ApiResponseCategory, bool)
	SetCachedCategoryCache(ctx context.Context, data *response.ApiResponseCategory)
}

type CategoryCommandCache interface {
	DeleteCachedCategoryCache(ctx context.Context, id int)
}

type CategoryStatsCache interface {
	GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) (*response.ApiResponseCategoryMonthlyTotalPrice, bool)
	SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data *response.ApiResponseCategoryMonthlyTotalPrice)

	GetCachedYearTotalPriceCache(ctx context.Context, year int) (*response.ApiResponseCategoryYearlyTotalPrice, bool)
	SetCachedYearTotalPriceCache(ctx context.Context, year int, data *response.ApiResponseCategoryYearlyTotalPrice)

	GetCachedMonthPriceCache(ctx context.Context, year int) (*response.ApiResponseCategoryMonthPrice, bool)
	SetCachedMonthPriceCache(ctx context.Context, year int, data *response.ApiResponseCategoryMonthPrice)

	GetCachedYearPriceCache(ctx context.Context, year int) (*response.ApiResponseCategoryYearPrice, bool)
	SetCachedYearPriceCache(ctx context.Context, year int, data *response.ApiResponseCategoryYearPrice)
}

type CategoryStatsByIdCache interface {
	GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) (*response.ApiResponseCategoryMonthlyTotalPrice, bool)
	SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data *response.ApiResponseCategoryMonthlyTotalPrice)

	GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) (*response.ApiResponseCategoryYearlyTotalPrice, bool)
	SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data *response.ApiResponseCategoryYearlyTotalPrice)

	GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) (*response.ApiResponseCategoryMonthPrice, bool)
	SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data *response.ApiResponseCategoryMonthPrice)

	GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) (*response.ApiResponseCategoryYearPrice, bool)
	SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data *response.ApiResponseCategoryYearPrice)
}

type CategoryStatsByMerchantCache interface {
	GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant) (*response.ApiResponseCategoryMonthlyTotalPrice, bool)
	SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant, data *response.ApiResponseCategoryMonthlyTotalPrice)

	GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant) (*response.ApiResponseCategoryYearlyTotalPrice, bool)
	SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant, data *response.ApiResponseCategoryYearlyTotalPrice)

	GetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant) (*response.ApiResponseCategoryMonthPrice, bool)
	SetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant, data *response.ApiResponseCategoryMonthPrice)

	GetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant) (*response.ApiResponseCategoryYearPrice, bool)
	SetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant, data *response.ApiResponseCategoryYearPrice)
}
