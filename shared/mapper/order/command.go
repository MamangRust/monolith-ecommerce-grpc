package orderapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type orderCommandResponseMapper struct{}

func NewOrderCommandResponseMapper() OrderCommandResponseMapper {
	return &orderCommandResponseMapper{}
}

func (o *orderCommandResponseMapper) ToResponseOrder(order *pb.OrderResponse) *response.OrderResponse {
	if order == nil { return nil }
	return &response.OrderResponse{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		UserID:     int(order.UserId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderCommandResponseMapper) ToResponsesOrder(orders []*pb.OrderResponse) []*response.OrderResponse {
	var mappedOrders []*response.OrderResponse
	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrder(order))
	}
	return mappedOrders
}

func (o *orderCommandResponseMapper) ToApiResponseOrder(pbResponse *pb.ApiResponseOrder) *response.ApiResponseOrder {
	return &response.ApiResponseOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrder(pbResponse.Data),
	}
}

func (o *orderCommandResponseMapper) ToResponseOrderDeleteAt(order *pb.OrderResponseDeleteAt) *response.OrderResponseDeleteAt {
	var deletedAt string
	if order.DeletedAt != nil {
		deletedAt = order.DeletedAt.Value
	}

	return &response.OrderResponseDeleteAt{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		UserID:     int(order.UserId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
		DeletedAt:  &deletedAt,
	}
}

func (o *orderCommandResponseMapper) ToResponsesOrderDeleteAt(orders []*pb.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt {
	var mappedOrders []*response.OrderResponseDeleteAt
	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrderDeleteAt(order))
	}
	return mappedOrders
}

func (o *orderCommandResponseMapper) ToApiResponseOrderDeleteAt(pbResponse *pb.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt {
	return &response.ApiResponseOrderDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderDeleteAt(pbResponse.Data),
	}
}

func (o *orderCommandResponseMapper) ToApiResponseOrderDelete(pbResponse *pb.ApiResponseOrderDelete) *response.ApiResponseOrderDelete {
	return &response.ApiResponseOrderDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderCommandResponseMapper) ToApiResponseOrderAll(pbResponse *pb.ApiResponseOrderAll) *response.ApiResponseOrderAll {
	return &response.ApiResponseOrderAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}
