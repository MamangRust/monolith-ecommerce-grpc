package orderitem_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItem, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItem)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItemDeleteAt, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItemDeleteAt)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) (*response.ApiResponsePaginationOrderItemDeleteAt, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data *response.ApiResponsePaginationOrderItemDeleteAt)

	GetCachedOrderItems(ctx context.Context, orderID int) (*response.ApiResponsesOrderItem, bool)
	SetCachedOrderItems(ctx context.Context, data *response.ApiResponsesOrderItem)
}

type OrderItemCommandCache interface {
	DeleteCachedOrderItemByOrderId(ctx context.Context, orderId int)
}
