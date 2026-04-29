package cartapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type cartQueryResponseMapper struct{}

func NewCartQueryResponseMapper() CartQueryResponseMapper {
	return &cartQueryResponseMapper{}
}

func (t *cartQueryResponseMapper) ToResponseCart(pbResponse *pb.CartResponse) *response.CartResponse {
	return &response.CartResponse{
		ID:        int(pbResponse.Id),
		UserID:    int(pbResponse.UserId),
		ProductID: int(pbResponse.ProductId),
		Name:      pbResponse.Name,
		Price:     int(pbResponse.Price),
		Image:     pbResponse.Image,
		Quantity:  int(pbResponse.Quantity),
		Weight:    int(pbResponse.Weight),
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
	}
}

func (t *cartQueryResponseMapper) ToResponseCarts(pbResponse []*pb.CartResponse) []*response.CartResponse {
	var carts []*response.CartResponse
	for _, cart := range pbResponse {
		carts = append(carts, t.ToResponseCart(cart))
	}
	return carts
}

func (t *cartQueryResponseMapper) ToApiResponseCart(pbResponse *pb.ApiResponseCart) *response.ApiResponseCart {
	return &response.ApiResponseCart{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    *t.ToResponseCart(pbResponse.Data),
	}
}

func (t *cartQueryResponseMapper) ToApiResponseCartPagination(pbResponse *pb.ApiResponsePaginationCart) *response.ApiResponseCartPagination {
	return &response.ApiResponseCartPagination{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       t.ToResponseCarts(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
