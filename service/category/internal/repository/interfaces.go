package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryStatsRepository interface {
	GetMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPrices(year int) ([]*record.CategoriesYearlyTotalPriceRecord, error)

	GetMonthPrice(year int) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPrice(year int) ([]*record.CategoriesYearPriceRecord, error)
}

type CategoryStatsByIdRepository interface {
	GetMonthlyTotalPriceById(req *requests.MonthTotalPriceCategory) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPricesById(req *requests.YearTotalPriceCategory) ([]*record.CategoriesYearlyTotalPriceRecord, error)

	GetMonthPriceById(req *requests.MonthPriceId) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPriceById(req *requests.YearPriceId) ([]*record.CategoriesYearPriceRecord, error)
}

type CategoryStatsByMerchantRepository interface {
	GetMonthlyTotalPriceByMerchant(req *requests.MonthTotalPriceMerchant) ([]*record.CategoriesMonthlyTotalPriceRecord, error)
	GetYearlyTotalPricesByMerchant(req *requests.YearTotalPriceMerchant) ([]*record.CategoriesYearlyTotalPriceRecord, error)

	GetMonthPriceByMerchant(req *requests.MonthPriceMerchant) ([]*record.CategoriesMonthPriceRecord, error)
	GetYearPriceByMerchant(req *requests.YearPriceMerchant) ([]*record.CategoriesYearPriceRecord, error)
}

type CategoryQueryRepository interface {
	FindAllCategory(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindByActive(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindByTrashed(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error)
	FindById(category_id int) (*record.CategoriesRecord, error)
	FindByIdTrashed(category_id int) (*record.CategoriesRecord, error)
}

type CategoryCommandRepository interface {
	CreateCategory(request *requests.CreateCategoryRequest) (*record.CategoriesRecord, error)
	UpdateCategory(request *requests.UpdateCategoryRequest) (*record.CategoriesRecord, error)
	TrashedCategory(category_id int) (*record.CategoriesRecord, error)
	RestoreCategory(category_id int) (*record.CategoriesRecord, error)
	DeleteCategoryPermanently(Category_id int) (bool, error)
	RestoreAllCategories() (bool, error)
	DeleteAllPermanentCategories() (bool, error)
}
