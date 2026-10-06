package productapimapper

import (
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type productCommandResponseMapper struct{}

func NewProductCommandResponseMapper() ProductCommandResponseMapper {
	return &productCommandResponseMapper{}
}

func (p *productCommandResponseMapper) ToResponseProduct(product *pbproduct.ProductResponse) *response.ProductResponse {
	if product == nil {
		return nil
	}
	return &response.ProductResponse{
		ID:           int(product.Id),
		MerchantID:   int(product.MerchantId),
		CategoryID:   int(product.CategoryId),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int(product.Weight),
		Rating:       float32(product.Rating),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func (p *productCommandResponseMapper) ToResponsesProduct(products []*pbproduct.ProductResponse) []*response.ProductResponse {
	var mappedProducts []*response.ProductResponse
	for _, product := range products {
		mappedProducts = append(mappedProducts, p.ToResponseProduct(product))
	}
	return mappedProducts
}

func (p *productCommandResponseMapper) ToResponseProductDeleteAt(product *pbproduct.ProductResponseDeleteAt) *response.ProductResponseDeleteAt {
	if product == nil {
		return nil
	}
	var deletedAt string
	if product.DeletedAt != nil {
		deletedAt = product.DeletedAt.Value
	}

	return &response.ProductResponseDeleteAt{
		ID:           int(product.Id),
		MerchantID:   int(product.MerchantId),
		CategoryID:   int(product.CategoryId),
		Name:         product.Name,
		Description:  product.Description,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand,
		Weight:       int(product.Weight),
		Rating:       float32(product.Rating),
		SlugProduct:  product.SlugProduct,
		ImageProduct: product.ImageProduct,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		DeletedAt:    &deletedAt,
	}
}

func (p *productCommandResponseMapper) ToResponsesProductDeleteAt(products []*pbproduct.ProductResponseDeleteAt) []*response.ProductResponseDeleteAt {
	var mappedProducts []*response.ProductResponseDeleteAt
	for _, product := range products {
		mappedProducts = append(mappedProducts, p.ToResponseProductDeleteAt(product))
	}
	return mappedProducts
}

func (p *productCommandResponseMapper) ToApiResponseProduct(pbResponse *pbproduct.ApiResponseProduct) *response.ApiResponseProduct {
	return &response.ApiResponseProduct{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponseProduct(pbResponse.Data),
	}
}

func (p *productCommandResponseMapper) ToApiResponsesProductDeleteAt(pbResponse *pbproduct.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt {
	return &response.ApiResponseProductDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponseProductDeleteAt(pbResponse.Data),
	}
}

func (p *productCommandResponseMapper) ToApiResponseProductDelete(pbResponse *pbproduct.ApiResponseProductDelete) *response.ApiResponseProductDelete {
	return &response.ApiResponseProductDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (p *productCommandResponseMapper) ToApiResponseProductAll(pbResponse *pbproduct.ApiResponseProductAll) *response.ApiResponseProductAll {
	return &response.ApiResponseProductAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (p *productCommandResponseMapper) ToApiResponsePaginationProductDeleteAt(pbResponse *pbproduct.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt {
	return nil // Not strictly needed in Command but to satisfy interface if we keep it there
}
