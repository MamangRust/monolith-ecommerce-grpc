package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*response.OrderItemResponse, total *int)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*response.OrderItemResponseDeleteAt, total *int)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*response.OrderItemResponseDeleteAt, total *int)

	GetCachedOrderItems(ctx context.Context, orderID int) ([]*response.OrderItemResponse, bool)
	SetCachedOrderItems(ctx context.Context, data []*response.OrderItemResponse)
}
