package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type orderItemRepository struct {
	client pb.OrderItemQueryServiceClient
}

func NewOrderItemRepository(client pb.OrderItemQueryServiceClient) *orderItemRepository {
	return &orderItemRepository{
		client: client,
	}
}

func (r *orderItemRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.GetOrderItemsByOrderRow, error) {
	res, err := r.client.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}

	var items []*db.GetOrderItemsByOrderRow
	for _, item := range res.Data {
		items = append(items, &db.GetOrderItemsByOrderRow{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    int32(item.Quantity),
			Price:       int32(item.Price),
		})
	}

	return items, nil
}
