package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryStatsRepository interface {
	GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPrices(ctx context.Context, year int) ([]*record.CategoriesYearlyTotalPriceRecord, error)

	GetMonthPrice(ctx context.Context, year int) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPrice(ctx context.Context, year int) ([]*record.CategoriesYearPriceRecord, error)
}

type CategoryStatsByIdRepository interface {
	GetMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*record.CategoriesYearlyTotalPriceRecord, error)

	GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*record.CategoriesYearPriceRecord, error)
}

type CategoryStatsByMerchantRepository interface {
	GetMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*record.CategoriesYearlyTotalPriceRecord, error)

	GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*record.CategoriesYearPriceRecord, error)
}

type CategoryQueryRepository interface {
	FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindById(ctx context.Context, categoryID int) (*record.CategoriesRecord, error)
	FindByIdTrashed(ctx context.Context, categoryID int) (*record.CategoriesRecord, error)
}

type CategoryCommandRepository interface {
	CreateCategory(ctx context.Context, request *requests.CreateCategoryRequest) (*record.CategoriesRecord, error)
	UpdateCategory(ctx context.Context, request *requests.UpdateCategoryRequest) (*record.CategoriesRecord, error)
	TrashedCategory(ctx context.Context, categoryID int) (*record.CategoriesRecord, error)
	RestoreCategory(ctx context.Context, categoryID int) (*record.CategoriesRecord, error)
	DeleteCategoryPermanently(ctx context.Context, categoryID int) (bool, error)
	RestoreAllCategories(ctx context.Context) (bool, error)
	DeleteAllPermanentCategories(ctx context.Context) (bool, error)
}
