package cartapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type cartCommandResponseMapper struct{}

func NewCartCommandResponseMapper() CartCommandResponseMapper {
	return &cartCommandResponseMapper{}
}

func (t *cartCommandResponseMapper) ToResponseCart(pbResponse *pb.CartResponse) *response.CartResponse {
	if pbResponse == nil { return nil }
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

func (t *cartCommandResponseMapper) ToResponseCarts(pbResponse []*pb.CartResponse) []*response.CartResponse {
	var carts []*response.CartResponse
	for _, cart := range pbResponse {
		carts = append(carts, t.ToResponseCart(cart))
	}
	return carts
}

func (t *cartCommandResponseMapper) ToApiResponseCart(pbResponse *pb.ApiResponseCart) *response.ApiResponseCart {
	return &response.ApiResponseCart{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    *t.ToResponseCart(pbResponse.Data),
	}
}

func (t *cartCommandResponseMapper) ToApiResponseCartDelete(pbResponse *pb.ApiResponseCartDelete) *response.ApiResponseCartDelete {
	return &response.ApiResponseCartDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (t *cartCommandResponseMapper) ToApiResponseCartAll(pbResponse *pb.ApiResponseCartAll) *response.ApiResponseCartAll {
	return &response.ApiResponseCartAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}
