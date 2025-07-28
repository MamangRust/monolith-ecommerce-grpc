package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CategoryStatsService interface {
	FindMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPrice(ctx context.Context, year int) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	FindMonthPrice(ctx context.Context, year int) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPrice(ctx context.Context, year int) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryStatsByIdService interface {
	FindMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPriceById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	FindMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryStatsByMerchantService interface {
	FindMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPriceByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	FindMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, categoryID int) (*response.CategoryResponse, *response.ErrorResponse)
}

type CategoryCommandService interface {
	CreateCategory(ctx context.Context, req *requests.CreateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse)
	UpdateCategory(ctx context.Context, req *requests.UpdateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse)
	TrashedCategory(ctx context.Context, categoryID int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse)
	RestoreCategory(ctx context.Context, categoryID int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse)
	DeleteCategoryPermanent(ctx context.Context, categoryID int) (bool, *response.ErrorResponse)
	RestoreAllCategories(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllCategoriesPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
