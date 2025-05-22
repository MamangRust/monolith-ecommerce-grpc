package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CategoryStatsService interface {
	FindMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPrice(year int) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	FindMonthPrice(year int) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPrice(year int) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryStatsByIdService interface {
	FindMonthlyTotalPriceById(req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPriceById(req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	FindMonthPriceById(req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPriceById(req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryStatsByMerchantService interface {
	FindMonthlyTotalPriceByMerchant(req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse)
	FindYearlyTotalPriceByMerchant(req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse)

	FindMonthPriceByMerchant(req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse)
	FindYearPriceByMerchant(req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse)
}

type CategoryQueryService interface {
	FindAll(req *requests.FindAllCategory) ([]*response.CategoryResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllCategory) ([]*response.CategoryResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(category_id int) (*response.CategoryResponse, *response.ErrorResponse)
}

type CategoryCommandService interface {
	CreateCategory(req *requests.CreateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse)
	UpdateCategory(req *requests.UpdateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse)
	TrashedCategory(category_id int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse)
	RestoreCategory(categoryID int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse)
	DeleteCategoryPermanent(categoryID int) (bool, *response.ErrorResponse)
	RestoreAllCategories() (bool, *response.ErrorResponse)
	DeleteAllCategoriesPermanent() (bool, *response.ErrorResponse)
}
