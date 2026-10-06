package categoryapimapper

import (
	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type categoryQueryResponseMapper struct {
	CategoryStatsResponseMapper
	CategoryCommandResponseMapper
}

func NewCategoryQueryResponseMapper() CategoryQueryResponseMapper {
	return &categoryQueryResponseMapper{
		CategoryStatsResponseMapper:   NewCategoryStatsResponseMapper(),
		CategoryCommandResponseMapper: NewCategoryCommandResponseMapper(),
	}
}

func (c *categoryQueryResponseMapper) ToResponseCategory(category *pbcategory.CategoryResponse) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID:            int(category.Id),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}
}

func (c *categoryQueryResponseMapper) ToResponsesCategory(categories []*pbcategory.CategoryResponse) []*response.CategoryResponse {
	var mappedCategories []*response.CategoryResponse
	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.ToResponseCategory(category))
	}
	return mappedCategories
}

func (c *categoryQueryResponseMapper) ToApiResponseCategory(pbResponse *pbcategory.ApiResponseCategory) *response.ApiResponseCategory {
	return &response.ApiResponseCategory{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategory(pbResponse.Data),
	}
}

func (c *categoryQueryResponseMapper) ToApiResponsesCategory(pbResponse *pbcategory.ApiResponsesCategory) *response.ApiResponsesCategory {
	return &response.ApiResponsesCategory{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponsesCategory(pbResponse.Data),
	}
}

func (c *categoryQueryResponseMapper) ToApiResponsePaginationCategory(pbResponse *pbcategory.ApiResponsePaginationCategory) *response.ApiResponsePaginationCategory {
	return &response.ApiResponsePaginationCategory{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       c.ToResponsesCategory(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
