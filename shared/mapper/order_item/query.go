package orderitemapimapper

import (
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type orderItemQueryResponseMapper struct{}

func NewOrderItemQueryResponseMapper() OrderItemQueryResponseMapper {
	return &orderItemQueryResponseMapper{}
}

func (o *orderItemQueryResponseMapper) ToResponseOrderItem(orderItem *pborder_item.OrderItemResponse) *response.OrderItemResponse {
	if orderItem == nil {
		return nil
	}
	return &response.OrderItemResponse{
		ID:        int(orderItem.Id),
		OrderID:   int(orderItem.OrderId),
		ProductID: int(orderItem.ProductId),
		Quantity:  int(orderItem.Quantity),
		Price:     int(orderItem.Price),
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
	}
}

func (o *orderItemQueryResponseMapper) ToResponsesOrderItem(orderItems []*pborder_item.OrderItemResponse) []*response.OrderItemResponse {
	var mappedOrderItems []*response.OrderItemResponse
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.ToResponseOrderItem(orderItem))
	}
	return mappedOrderItems
}

func (o *orderItemQueryResponseMapper) ToApiResponseOrderItem(pbResponse *pborder_item.ApiResponseOrderItem) *response.ApiResponseOrderItem {
	return &response.ApiResponseOrderItem{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderItem(pbResponse.Data),
	}
}

func (o *orderItemQueryResponseMapper) ToApiResponsesOrderItem(pbResponse *pborder_item.ApiResponsesOrderItem) *response.ApiResponsesOrderItem {
	return &response.ApiResponsesOrderItem{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponsesOrderItem(pbResponse.Data),
	}
}

func (o *orderItemQueryResponseMapper) ToApiResponsePaginationOrderItem(pbResponse *pborder_item.ApiResponsePaginationOrderItem) *response.ApiResponsePaginationOrderItem {
	return &response.ApiResponsePaginationOrderItem{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrderItem(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (o *orderItemQueryResponseMapper) ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pborder_item.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt {
	var mappedOrderItems []*response.OrderItemResponseDeleteAt
	for _, orderItem := range pbResponse.Data {
		var deletedAt string
		if orderItem.DeletedAt != nil {
			deletedAt = orderItem.DeletedAt.Value
		}
		mappedOrderItems = append(mappedOrderItems, &response.OrderItemResponseDeleteAt{
			ID:        int(orderItem.Id),
			OrderID:   int(orderItem.OrderId),
			ProductID: int(orderItem.ProductId),
			Quantity:  int(orderItem.Quantity),
			Price:     int(orderItem.Price),
			CreatedAt: orderItem.CreatedAt,
			UpdatedAt: orderItem.UpdatedAt,
			DeletedAt: &deletedAt,
		})
	}
	return &response.ApiResponsePaginationOrderItemDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       mappedOrderItems,
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
