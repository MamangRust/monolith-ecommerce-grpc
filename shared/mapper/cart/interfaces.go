package cartapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CartBaseResponseMapper interface {
	ToResponseCart(pbResponse *pb.CartResponse) *response.CartResponse
	ToResponseCarts(pbResponse []*pb.CartResponse) []*response.CartResponse
	ToApiResponseCart(pbResponse *pb.ApiResponseCart) *response.ApiResponseCart
}

type CartQueryResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartPagination(pbResponse *pb.ApiResponsePaginationCart) *response.ApiResponseCartPagination
}

type CartCommandResponseMapper interface {
	CartBaseResponseMapper
	ToApiResponseCartDelete(pbResponse *pb.ApiResponseCartDelete) *response.ApiResponseCartDelete
	ToApiResponseCartAll(pbResponse *pb.ApiResponseCartAll) *response.ApiResponseCartAll
}
