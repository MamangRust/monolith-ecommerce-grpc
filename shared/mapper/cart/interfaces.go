package cartapimapper

import (
	pbcart "github.com/MamangRust/monolith-ecommerce-pb/cart"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CartBaseResponseMapper interface {
	ToResponseCart(pbResponse *pbcart.CartResponse) *response.CartResponse
	ToResponseCarts(pbResponse []*pbcart.CartResponse) []*response.CartResponse
	ToApiResponseCart(pbResponse *pbcart.ApiResponseCart) *response.ApiResponseCart
}

type CartQueryResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartPagination(pbResponse *pbcart.ApiResponsePaginationCart) *response.ApiResponseCartPagination
}

type CartCommandResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartDelete(pbResponse *pbcart.ApiResponseCartDelete) *response.ApiResponseCartDelete
	ToApiResponseCartAll(pbResponse *pbcart.ApiResponseCartAll) *response.ApiResponseCartAll
}
