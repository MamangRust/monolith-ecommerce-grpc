package orderitemapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderItemBaseResponseMapper interface {
	ToResponseOrderItem(orderItem *pb.OrderItemResponse) *response.OrderItemResponse
	ToResponsesOrderItem(orderItems []*pb.OrderItemResponse) []*response.OrderItemResponse
}

type OrderItemQueryResponseMapper interface {
	OrderItemBaseResponseMapper
	ToApiResponseOrderItem(pbResponse *pb.ApiResponseOrderItem) *response.ApiResponseOrderItem
	ToApiResponsesOrderItem(pbResponse *pb.ApiResponsesOrderItem) *response.ApiResponsesOrderItem
	ToApiResponsePaginationOrderItem(pbResponse *pb.ApiResponsePaginationOrderItem) *response.ApiResponsePaginationOrderItem
	ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pb.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt
}

type OrderItemCommandResponseMapper interface {
	OrderItemBaseResponseMapper
	ToApiResponseOrderItem(pbResponse *pb.ApiResponseOrderItem) *response.ApiResponseOrderItem
	ToResponseOrderItemDeleteAt(orderItem *pb.OrderItemResponseDeleteAt) *response.OrderItemResponseDeleteAt
	ToResponsesOrderItemDeleteAt(orderItems []*pb.OrderItemResponseDeleteAt) []*response.OrderItemResponseDeleteAt
	ToApiResponseOrderItemDelete(pbResponse *pb.ApiResponseOrderItemDelete) *response.ApiResponseOrderItemDelete
	ToApiResponseOrderItemAll(pbResponse *pb.ApiResponseOrderItemAll) *response.ApiResponseOrderItemAll
	ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pb.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt
}
