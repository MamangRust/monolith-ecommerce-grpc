package productapimapper

import (
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ProductBaseResponseMapper interface {
	ToResponseProduct(product *pbproduct.ProductResponse) *response.ProductResponse
	ToResponsesProduct(products []*pbproduct.ProductResponse) []*response.ProductResponse
	ToResponseProductDeleteAt(product *pbproduct.ProductResponseDeleteAt) *response.ProductResponseDeleteAt
	ToResponsesProductDeleteAt(products []*pbproduct.ProductResponseDeleteAt) []*response.ProductResponseDeleteAt
	ToApiResponseProduct(pbResponse *pbproduct.ApiResponseProduct) *response.ApiResponseProduct
	ToApiResponsePaginationProductDeleteAt(pbResponse *pbproduct.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt
}

type ProductQueryResponseMapper interface {
	ProductBaseResponseMapper
	ToApiResponsesProduct(pbResponse *pbproduct.ApiResponsesProduct) *response.ApiResponsesProduct
	ToApiResponsePaginationProduct(pbResponse *pbproduct.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct
}

type ProductCommandResponseMapper interface {
	ProductBaseResponseMapper
	ToApiResponsesProductDeleteAt(pbResponse *pbproduct.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt
	ToApiResponseProductDelete(pbResponse *pbproduct.ApiResponseProductDelete) *response.ApiResponseProductDelete
	ToApiResponseProductAll(pbResponse *pbproduct.ApiResponseProductAll) *response.ApiResponseProductAll
}
