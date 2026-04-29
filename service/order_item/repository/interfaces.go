package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderItemQueryRepository interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsRow, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsActiveRow, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsTrashedRow, error)

	FindOrderItemByOrder(
		ctx context.Context,
		order_id int,
	) ([]*db.GetOrderItemsByOrderRow, error)
}

type OrderItemCommandRepository interface {
	Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error)
	Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error)

	Trash(ctx context.Context, orderItemID int) (*db.OrderItem, error)
	Restore(ctx context.Context, orderItemID int) (*db.OrderItem, error)
	DeletePermanent(ctx context.Context, orderItemID int) (bool, error)
	DeleteOrderItemByOrderPermanent(ctx context.Context, orderID int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)

	CalculateTotalPrice(ctx context.Context, orderID int) (int, error)
}
