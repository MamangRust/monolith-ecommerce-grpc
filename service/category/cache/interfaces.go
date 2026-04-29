package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, *int, bool)
	SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*db.GetCategoriesRow, total *int)

	GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, *int, bool)
	SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*db.GetCategoriesActiveRow, total *int)

	GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, *int, bool)
	SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*db.GetCategoriesTrashedRow, total *int)

	GetCachedCategoryCache(ctx context.Context, id int) (*db.GetCategoryByIDRow, bool)
	SetCachedCategoryCache(ctx context.Context, data *db.GetCategoryByIDRow)
}

type CategoryCommandCache interface {
	DeleteCachedCategoryCache(ctx context.Context, id int)
}

type CategoryStatsCache interface {
	GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, bool)
	SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data []*db.GetMonthlyTotalPriceRow)

	GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, bool)
	SetCachedYearTotalPriceCache(ctx context.Context, year int, data []*db.GetYearlyTotalPriceRow)

	GetCachedMonthPriceCache(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, bool)
	SetCachedMonthPriceCache(ctx context.Context, year int, data []*db.GetMonthlyCategoryRow)

	GetCachedYearPriceCache(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, bool)
	SetCachedYearPriceCache(ctx context.Context, year int, data []*db.GetYearlyCategoryRow)
}

type CategoryStatsByIdCache interface {
	GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*db.GetMonthlyTotalPriceByIdRow, bool)
	SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data []*db.GetMonthlyTotalPriceByIdRow)

	GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*db.GetYearlyTotalPriceByIdRow, bool)
	SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data []*db.GetYearlyTotalPriceByIdRow)

	GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) ([]*db.GetMonthlyCategoryByIdRow, bool)
	SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data []*db.GetMonthlyCategoryByIdRow)

	GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) ([]*db.GetYearlyCategoryByIdRow, bool)
	SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data []*db.GetYearlyCategoryByIdRow)
}

type CategoryStatsByMerchantCache interface {
	GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*db.GetMonthlyTotalPriceByMerchantRow, bool)
	SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant, data []*db.GetMonthlyTotalPriceByMerchantRow)

	GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*db.GetYearlyTotalPriceByMerchantRow, bool)
	SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant, data []*db.GetYearlyTotalPriceByMerchantRow)

	GetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant) ([]*db.GetMonthlyCategoryByMerchantRow, bool)
	SetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant, data []*db.GetMonthlyCategoryByMerchantRow)

	GetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant) ([]*db.GetYearlyCategoryByMerchantRow, bool)
	SetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant, data []*db.GetYearlyCategoryByMerchantRow)
}
