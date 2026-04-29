package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryStatsRepository interface {
	GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, error)
	GetYearlyTotalPrices(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, error)
	GetMonthPrice(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, error)
	GetYearPrice(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, error)
}

type CategoryStatsByIdRepository interface {
	GetMonthlyTotalPriceById(
		ctx context.Context,
		req *requests.MonthTotalPriceCategory,
	) ([]*db.GetMonthlyTotalPriceByIdRow, error)
	GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*db.GetYearlyTotalPriceByIdRow, error)
	GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*db.GetMonthlyCategoryByIdRow, error)
	GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*db.GetYearlyCategoryByIdRow, error)
}

type CategoryStatsByMerchantRepository interface {
	GetMonthlyTotalPriceByMerchant(
		ctx context.Context,
		req *requests.MonthTotalPriceMerchant,
	) ([]*db.GetMonthlyTotalPriceByMerchantRow, error)
	GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*db.GetYearlyTotalPriceByMerchantRow, error)
	GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*db.GetMonthlyCategoryByMerchantRow, error)
	GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*db.GetYearlyCategoryByMerchantRow, error)
}

type CategoryQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, error)

	FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, error)

	FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, error)

	FindByID(ctx context.Context, category_id int) (*db.GetCategoryByIDRow, error)
	FindByIDTrashed(ctx context.Context, category_id int) (*db.Category, error)
}

type CategoryCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateCategoryRequest,
	) (*db.CreateCategoryRow, error)

	Update(
		ctx context.Context,
		request *requests.UpdateCategoryRequest,
	) (*db.UpdateCategoryRow, error)

	Trash(
		ctx context.Context,
		category_id int,
	) (*db.Category, error)

	Restore(
		ctx context.Context,
		category_id int,
	) (*db.Category, error)

	DeletePermanent(
		ctx context.Context,
		category_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
