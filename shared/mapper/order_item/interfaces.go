package orderitemapimapper

import (
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderItemBaseResponseMapper interface {
	ToResponseOrderItem(orderItem *pborder_item.OrderItemResponse) *response.OrderItemResponse
	ToResponsesOrderItem(orderItems []*pborder_item.OrderItemResponse) []*response.OrderItemResponse
}

type OrderItemQueryResponseMapper interface {
	OrderItemBaseResponseMapper
	ToApiResponseOrderItem(pbResponse *pborder_item.ApiResponseOrderItem) *response.ApiResponseOrderItem
	ToApiResponsesOrderItem(pbResponse *pborder_item.ApiResponsesOrderItem) *response.ApiResponsesOrderItem
	ToApiResponsePaginationOrderItem(pbResponse *pborder_item.ApiResponsePaginationOrderItem) *response.ApiResponsePaginationOrderItem
	ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pborder_item.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt
}

type OrderItemCommandResponseMapper interface {
	OrderItemBaseResponseMapper
	ToApiResponseOrderItem(pbResponse *pborder_item.ApiResponseOrderItem) *response.ApiResponseOrderItem
	ToResponseOrderItemDeleteAt(orderItem *pborder_item.OrderItemResponseDeleteAt) *response.OrderItemResponseDeleteAt
	ToResponsesOrderItemDeleteAt(orderItems []*pborder_item.OrderItemResponseDeleteAt) []*response.OrderItemResponseDeleteAt
	ToApiResponseOrderItemDelete(pbResponse *pborder_item.ApiResponseOrderItemDelete) *response.ApiResponseOrderItemDelete
	ToApiResponseOrderItemAll(pbResponse *pborder_item.ApiResponseOrderItemAll) *response.ApiResponseOrderItemAll
	ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pborder_item.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt
}
