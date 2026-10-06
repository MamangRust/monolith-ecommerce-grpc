package productapimapper

import (
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type productQueryResponseMapper struct{}

func NewProductQueryResponseMapper() ProductQueryResponseMapper {
	return &productQueryResponseMapper{}
}

func (p *productQueryResponseMapper) ToResponseProduct(product *pbproduct.ProductResponse) *response.ProductResponse {
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

func (p *productQueryResponseMapper) ToResponsesProduct(products []*pbproduct.ProductResponse) []*response.ProductResponse {
	var mappedProducts []*response.ProductResponse
	for _, product := range products {
		mappedProducts = append(mappedProducts, p.ToResponseProduct(product))
	}
	return mappedProducts
}

func (p *productQueryResponseMapper) ToResponseProductDeleteAt(product *pbproduct.ProductResponseDeleteAt) *response.ProductResponseDeleteAt {
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

func (p *productQueryResponseMapper) ToResponsesProductDeleteAt(products []*pbproduct.ProductResponseDeleteAt) []*response.ProductResponseDeleteAt {
	var mappedProducts []*response.ProductResponseDeleteAt
	for _, product := range products {
		mappedProducts = append(mappedProducts, p.ToResponseProductDeleteAt(product))
	}
	return mappedProducts
}

func (p *productQueryResponseMapper) ToApiResponseProduct(pbResponse *pbproduct.ApiResponseProduct) *response.ApiResponseProduct {
	return &response.ApiResponseProduct{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponseProduct(pbResponse.Data),
	}
}

func (p *productQueryResponseMapper) ToApiResponsesProduct(pbResponse *pbproduct.ApiResponsesProduct) *response.ApiResponsesProduct {
	return &response.ApiResponsesProduct{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    p.ToResponsesProduct(pbResponse.Data),
	}
}

func (p *productQueryResponseMapper) ToApiResponsePaginationProduct(pbResponse *pbproduct.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct {
	return &response.ApiResponsePaginationProduct{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       p.ToResponsesProduct(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (p *productQueryResponseMapper) ToApiResponsePaginationProductDeleteAt(pbResponse *pbproduct.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt {
	return &response.ApiResponsePaginationProductDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       p.ToResponsesProductDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
