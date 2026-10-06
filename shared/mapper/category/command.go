package categoryapimapper

import (
	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type categoryCommandResponseMapper struct{}

func NewCategoryCommandResponseMapper() CategoryCommandResponseMapper {
	return &categoryCommandResponseMapper{}
}

func (c *categoryCommandResponseMapper) ToResponseCategory(category *pbcategory.CategoryResponse) *response.CategoryResponse {
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

func (c *categoryCommandResponseMapper) ToResponsesCategory(categories []*pbcategory.CategoryResponse) []*response.CategoryResponse {
	var mappedCategories []*response.CategoryResponse
	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.ToResponseCategory(category))
	}
	return mappedCategories
}

func (c *categoryCommandResponseMapper) ToResponseCategoryDelete(category *pbcategory.CategoryResponseDeleteAt) *response.CategoryResponseDeleteAt {
	var deletedAt string
	if category.DeletedAt != nil {
		deletedAt = category.DeletedAt.Value
	}

	return &response.CategoryResponseDeleteAt{
		ID:            int(category.Id),
		Name:          category.Name,
		Description:   category.Description,
		SlugCategory:  category.SlugCategory,
		ImageCategory: category.ImageCategory,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
		DeletedAt:     &deletedAt,
	}
}

func (c *categoryCommandResponseMapper) ToResponsesCategoryDeleteAt(categories []*pbcategory.CategoryResponseDeleteAt) []*response.CategoryResponseDeleteAt {
	var mappedCategories []*response.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, c.ToResponseCategoryDelete(category))
	}
	return mappedCategories
}

func (c *categoryCommandResponseMapper) ToApiResponseCategory(pbResponse *pbcategory.ApiResponseCategory) *response.ApiResponseCategory {
	return &response.ApiResponseCategory{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategory(pbResponse.Data),
	}
}

func (c *categoryCommandResponseMapper) ToApiResponseCategoryDeleteAt(pbResponse *pbcategory.ApiResponseCategoryDeleteAt) *response.ApiResponseCategoryDeleteAt {
	return &response.ApiResponseCategoryDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    c.ToResponseCategoryDelete(pbResponse.Data),
	}
}

func (c *categoryCommandResponseMapper) ToApiResponseCategoryDelete(pbResponse *pbcategory.ApiResponseCategoryDelete) *response.ApiResponseCategoryDelete {
	return &response.ApiResponseCategoryDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (c *categoryCommandResponseMapper) ToApiResponseCategoryAll(pbResponse *pbcategory.ApiResponseCategoryAll) *response.ApiResponseCategoryAll {
	return &response.ApiResponseCategoryAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (c *categoryCommandResponseMapper) ToApiResponsePaginationCategoryDeleteAt(pbResponse *pbcategory.ApiResponsePaginationCategoryDeleteAt) *response.ApiResponsePaginationCategoryDeleteAt {
	return &response.ApiResponsePaginationCategoryDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       c.ToResponsesCategoryDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
