package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderItemQueryService interface {
	FindAllOrderItems(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsTrashedRow, *int, error)

	FindOrderItemByOrder(
		ctx context.Context,
		order_id int,
	) ([]*db.GetOrderItemsByOrderRow, error)
}

type OrderItemCommandService interface {
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error)

	TrashOrderItem(ctx context.Context, orderItemID int) (*db.OrderItem, error)
	RestoreOrderItem(ctx context.Context, orderItemID int) (*db.OrderItem, error)
	DeleteOrderItemPermanent(ctx context.Context, orderItemID int) (bool, error)

	RestoreAllOrdersItem(ctx context.Context) (bool, error)
	DeleteAllPermanentOrdersItem(ctx context.Context) (bool, error)

	CalculateTotalPrice(ctx context.Context, orderID int) (int, error)
}
