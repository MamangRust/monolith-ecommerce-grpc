package orderapimapper

import (
	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type orderCommandResponseMapper struct{}

func NewOrderCommandResponseMapper() OrderCommandResponseMapper {
	return &orderCommandResponseMapper{}
}

func (o *orderCommandResponseMapper) ToResponseOrder(order *pborder.OrderResponse) *response.OrderResponse {
	if order == nil {
		return nil
	}
	return &response.OrderResponse{
		ID:         int(order.Id),
		MerchantID: int(order.MerchantId),
		UserID:     int(order.UserId),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

func (o *orderCommandResponseMapper) ToResponsesOrder(orders []*pborder.OrderResponse) []*response.OrderResponse {
	var mappedOrders []*response.OrderResponse
	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrder(order))
	}
	return mappedOrders
}

func (o *orderCommandResponseMapper) ToApiResponseOrder(pbResponse *pborder.ApiResponseOrder) *response.ApiResponseOrder {
	return &response.ApiResponseOrder{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrder(pbResponse.Data),
	}
}

func (o *orderCommandResponseMapper) ToResponseOrderDeleteAt(order *pborder.OrderResponseDeleteAt) *response.OrderResponseDeleteAt {
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

func (o *orderCommandResponseMapper) ToResponsesOrderDeleteAt(orders []*pborder.OrderResponseDeleteAt) []*response.OrderResponseDeleteAt {
	var mappedOrders []*response.OrderResponseDeleteAt
	for _, order := range orders {
		mappedOrders = append(mappedOrders, o.ToResponseOrderDeleteAt(order))
	}
	return mappedOrders
}

func (o *orderCommandResponseMapper) ToApiResponseOrderDeleteAt(pbResponse *pborder.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt {
	return &response.ApiResponseOrderDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderDeleteAt(pbResponse.Data),
	}
}

func (o *orderCommandResponseMapper) ToApiResponseOrderDelete(pbResponse *pborder.ApiResponseOrderDelete) *response.ApiResponseOrderDelete {
	return &response.ApiResponseOrderDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderCommandResponseMapper) ToApiResponseOrderAll(pbResponse *pborder.ApiResponseOrderAll) *response.ApiResponseOrderAll {
	return &response.ApiResponseOrderAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}
