package orderitemapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type orderItemCommandResponseMapper struct{}

func NewOrderItemCommandResponseMapper() OrderItemCommandResponseMapper {
	return &orderItemCommandResponseMapper{}
}

func (o *orderItemCommandResponseMapper) ToResponseOrderItem(orderItem *pb.OrderItemResponse) *response.OrderItemResponse {
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

func (o *orderItemCommandResponseMapper) ToResponsesOrderItem(orderItems []*pb.OrderItemResponse) []*response.OrderItemResponse {
	var mappedOrderItems []*response.OrderItemResponse
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.ToResponseOrderItem(orderItem))
	}
	return mappedOrderItems
}

func (o *orderItemCommandResponseMapper) ToApiResponseOrderItem(pbResponse *pb.ApiResponseOrderItem) *response.ApiResponseOrderItem {
	return &response.ApiResponseOrderItem{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    o.ToResponseOrderItem(pbResponse.Data),
	}
}

func (o *orderItemCommandResponseMapper) ToResponseOrderItemDeleteAt(orderItem *pb.OrderItemResponseDeleteAt) *response.OrderItemResponseDeleteAt {
	var deletedAt string
	if orderItem.DeletedAt != nil {
		deletedAt = orderItem.DeletedAt.Value
	}

	return &response.OrderItemResponseDeleteAt{
		ID:        int(orderItem.Id),
		OrderID:   int(orderItem.OrderId),
		ProductID: int(orderItem.ProductId),
		Quantity:  int(orderItem.Quantity),
		Price:     int(orderItem.Price),
		CreatedAt: orderItem.CreatedAt,
		UpdatedAt: orderItem.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (o *orderItemCommandResponseMapper) ToResponsesOrderItemDeleteAt(orderItems []*pb.OrderItemResponseDeleteAt) []*response.OrderItemResponseDeleteAt {
	var mappedOrderItems []*response.OrderItemResponseDeleteAt
	for _, orderItem := range orderItems {
		mappedOrderItems = append(mappedOrderItems, o.ToResponseOrderItemDeleteAt(orderItem))
	}
	return mappedOrderItems
}

func (o *orderItemCommandResponseMapper) ToApiResponseOrderItemDelete(pbResponse *pb.ApiResponseOrderItemDelete) *response.ApiResponseOrderItemDelete {
	return &response.ApiResponseOrderItemDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderItemCommandResponseMapper) ToApiResponseOrderItemAll(pbResponse *pb.ApiResponseOrderItemAll) *response.ApiResponseOrderItemAll {
	return &response.ApiResponseOrderItemAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (o *orderItemCommandResponseMapper) ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pb.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt {
	return &response.ApiResponsePaginationOrderItemDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       o.ToResponsesOrderItemDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
