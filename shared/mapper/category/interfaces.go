package categoryapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CategoryBaseResponseMapper interface {
	ToResponseCategory(category *pb.CategoryResponse) *response.CategoryResponse
	ToResponsesCategory(categories []*pb.CategoryResponse) []*response.CategoryResponse
}

type CategoryQueryResponseMapper interface {
	CategoryBaseResponseMapper
	CategoryStatsResponseMapper
	ToApiResponseCategory(pbResponse *pb.ApiResponseCategory) *response.ApiResponseCategory
	ToApiResponsesCategory(pbResponse *pb.ApiResponsesCategory) *response.ApiResponsesCategory
	ToApiResponsePaginationCategory(pbResponse *pb.ApiResponsePaginationCategory) *response.ApiResponsePaginationCategory
	ToApiResponsePaginationCategoryDeleteAt(pbResponse *pb.ApiResponsePaginationCategoryDeleteAt) *response.ApiResponsePaginationCategoryDeleteAt
}

type CategoryCommandResponseMapper interface {
	CategoryBaseResponseMapper
	ToResponseCategoryDelete(category *pb.CategoryResponseDeleteAt) *response.CategoryResponseDeleteAt
	ToResponsesCategoryDeleteAt(categories []*pb.CategoryResponseDeleteAt) []*response.CategoryResponseDeleteAt
	ToApiResponseCategoryDeleteAt(pbResponse *pb.ApiResponseCategoryDeleteAt) *response.ApiResponseCategoryDeleteAt
	ToApiResponseCategory(pbResponse *pb.ApiResponseCategory) *response.ApiResponseCategory
	ToApiResponseCategoryDelete(pbResponse *pb.ApiResponseCategoryDelete) *response.ApiResponseCategoryDelete
	ToApiResponseCategoryAll(pbResponse *pb.ApiResponseCategoryAll) *response.ApiResponseCategoryAll
	ToApiResponsePaginationCategoryDeleteAt(pbResponse *pb.ApiResponsePaginationCategoryDeleteAt) *response.ApiResponsePaginationCategoryDeleteAt
}

type CategoryStatsResponseMapper interface {
	ToResponseCategoryMonthlyPrice(category *pb.CategoryMonthPriceResponse) *response.CategoryMonthPriceResponse
	ToResponseCategoryMonthlyPrices(c []*pb.CategoryMonthPriceResponse) []*response.CategoryMonthPriceResponse
	ToResponseCategoryYearlyPrice(category *pb.CategoryYearPriceResponse) *response.CategoryYearPriceResponse
	ToResponseCategoryYearlyPrices(c []*pb.CategoryYearPriceResponse) []*response.CategoryYearPriceResponse
	ToResponseCashierMonthlyTotalPrice(c *pb.CategoriesMonthlyTotalPriceResponse) *response.CategoriesMonthlyTotalPriceResponse
	ToResponseCategoryMonthlyTotalPrices(c []*pb.CategoriesMonthlyTotalPriceResponse) []*response.CategoriesMonthlyTotalPriceResponse
	ToResponseCategoryYearlyTotalSale(c *pb.CategoriesYearlyTotalPriceResponse) *response.CategoriesYearlyTotalPriceResponse
	ToResponseCategoryYearlyTotalPrices(c []*pb.CategoriesYearlyTotalPriceResponse) []*response.CategoriesYearlyTotalPriceResponse

	ToApiResponseCategoryMonthPrice(pbResponse *pb.ApiResponseCategoryMonthPrice) *response.ApiResponseCategoryMonthPrice
	ToApiResponseCategoryYearPrice(pbResponse *pb.ApiResponseCategoryYearPrice) *response.ApiResponseCategoryYearPrice
	ToApiResponseCategoryMonthlyTotalPrice(pbResponse *pb.ApiResponseCategoryMonthlyTotalPrice) *response.ApiResponseCategoryMonthlyTotalPrice
	ToApiResponseCategoryYearlyTotalPrice(pbResponse *pb.ApiResponseCategoryYearlyTotalPrice) *response.ApiResponseCategoryYearlyTotalPrice
}
