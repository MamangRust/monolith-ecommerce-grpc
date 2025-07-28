package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type OrderItemQueryRepository interface {
	FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*record.OrderItemRecord, error)
}
