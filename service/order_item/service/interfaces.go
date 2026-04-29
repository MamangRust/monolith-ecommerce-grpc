package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderItemQueryService interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsRow, *int, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsActiveRow, *int, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllOrderItems,
	) ([]*db.GetOrderItemsTrashedRow, *int, error)

	FindByOrder(
		ctx context.Context,
		order_id int,
	) ([]*db.GetOrderItemsByOrderRow, error)
}

type OrderItemCommandService interface {
	Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error)
	Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error)

	Trash(ctx context.Context, orderItemID int) (*db.OrderItem, error)
	Restore(ctx context.Context, orderItemID int) (*db.OrderItem, error)
	DeletePermanent(ctx context.Context, orderItemID int) (bool, error)
	DeleteByOrderPermanent(ctx context.Context, orderID int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)

	CalculateTotalPrice(ctx context.Context, orderID int) (int, error)
}
