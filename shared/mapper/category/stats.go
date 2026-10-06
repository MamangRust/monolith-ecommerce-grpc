package categoryapimapper

import (
	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type categoryStatsResponseMapper struct{}

func NewCategoryStatsResponseMapper() CategoryStatsResponseMapper {
	return &categoryStatsResponseMapper{}
}

func (m *categoryStatsResponseMapper) ToResponseCategoryMonthlyPrice(category *pbcategory.CategoryMonthPriceResponse) *response.CategoryMonthPriceResponse {
	return &response.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryID:   int(category.CategoryId),
		CategoryName: category.CategoryName,
		OrderCount:   int(category.OrderCount),
		ItemsSold:    int(category.ItemsSold),
		TotalRevenue: int(category.TotalRevenue),
	}
}

func (m *categoryStatsResponseMapper) ToResponseCategoryMonthlyPrices(c []*pbcategory.CategoryMonthPriceResponse) []*response.CategoryMonthPriceResponse {
	var mapped []*response.CategoryMonthPriceResponse
	for _, item := range c {
		mapped = append(mapped, m.ToResponseCategoryMonthlyPrice(item))
	}
	return mapped
}

func (m *categoryStatsResponseMapper) ToResponseCategoryYearlyPrice(category *pbcategory.CategoryYearPriceResponse) *response.CategoryYearPriceResponse {
	return &response.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryID:         int(category.CategoryId),
		CategoryName:       category.CategoryName,
		OrderCount:         int(category.OrderCount),
		ItemsSold:          int(category.ItemsSold),
		TotalRevenue:       int(category.TotalRevenue),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (m *categoryStatsResponseMapper) ToResponseCategoryYearlyPrices(c []*pbcategory.CategoryYearPriceResponse) []*response.CategoryYearPriceResponse {
	var mapped []*response.CategoryYearPriceResponse
	for _, item := range c {
		mapped = append(mapped, m.ToResponseCategoryYearlyPrice(item))
	}
	return mapped
}

func (m *categoryStatsResponseMapper) ToResponseCashierMonthlyTotalPrice(c *pbcategory.CategoriesMonthlyTotalPriceResponse) *response.CategoriesMonthlyTotalPriceResponse {
	return &response.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (m *categoryStatsResponseMapper) ToResponseCategoryMonthlyTotalPrices(c []*pbcategory.CategoriesMonthlyTotalPriceResponse) []*response.CategoriesMonthlyTotalPriceResponse {
	var mapped []*response.CategoriesMonthlyTotalPriceResponse
	for _, item := range c {
		mapped = append(mapped, m.ToResponseCashierMonthlyTotalPrice(item))
	}
	return mapped
}

func (m *categoryStatsResponseMapper) ToResponseCategoryYearlyTotalSale(c *pbcategory.CategoriesYearlyTotalPriceResponse) *response.CategoriesYearlyTotalPriceResponse {
	return &response.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (m *categoryStatsResponseMapper) ToResponseCategoryYearlyTotalPrices(c []*pbcategory.CategoriesYearlyTotalPriceResponse) []*response.CategoriesYearlyTotalPriceResponse {
	var mapped []*response.CategoriesYearlyTotalPriceResponse
	for _, item := range c {
		mapped = append(mapped, m.ToResponseCategoryYearlyTotalSale(item))
	}
	return mapped
}

func (m *categoryStatsResponseMapper) ToApiResponseCategoryMonthPrice(pbResponse *pbcategory.ApiResponseCategoryMonthPrice) *response.ApiResponseCategoryMonthPrice {
	return &response.ApiResponseCategoryMonthPrice{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseCategoryMonthlyPrices(pbResponse.Data),
	}
}

func (m *categoryStatsResponseMapper) ToApiResponseCategoryYearPrice(pbResponse *pbcategory.ApiResponseCategoryYearPrice) *response.ApiResponseCategoryYearPrice {
	return &response.ApiResponseCategoryYearPrice{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseCategoryYearlyPrices(pbResponse.Data),
	}
}

func (m *categoryStatsResponseMapper) ToApiResponseCategoryMonthlyTotalPrice(pbResponse *pbcategory.ApiResponseCategoryMonthlyTotalPrice) *response.ApiResponseCategoryMonthlyTotalPrice {
	return &response.ApiResponseCategoryMonthlyTotalPrice{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseCategoryMonthlyTotalPrices(pbResponse.Data),
	}
}

func (m *categoryStatsResponseMapper) ToApiResponseCategoryYearlyTotalPrice(pbResponse *pbcategory.ApiResponseCategoryYearlyTotalPrice) *response.ApiResponseCategoryYearlyTotalPrice {
	return &response.ApiResponseCategoryYearlyTotalPrice{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseCategoryYearlyTotalPrices(pbResponse.Data),
	}
}
