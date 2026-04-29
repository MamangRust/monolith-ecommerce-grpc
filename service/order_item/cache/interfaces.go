package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, *int, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*db.GetOrderItemsRow, total *int)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, *int, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*db.GetOrderItemsActiveRow, total *int)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, *int, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*db.GetOrderItemsTrashedRow, total *int)

	GetCachedOrderItems(ctx context.Context, orderID int) ([]*db.GetOrderItemsByOrderRow, bool)
	SetCachedOrderItems(ctx context.Context, data []*db.GetOrderItemsByOrderRow)
}

type OrderItemCommandCache interface {
	InvalidateOrderItemCache(ctx context.Context) error
}
