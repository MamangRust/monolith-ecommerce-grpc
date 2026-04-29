package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	order_item_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type orderItemQueryRepository struct {
	queryClient   pb.OrderItemQueryServiceClient
	commandClient pb.OrderItemCommandServiceClient
}

func NewOrderItemQueryRepository(queryClient pb.OrderItemQueryServiceClient, commandClient pb.OrderItemCommandServiceClient) *orderItemQueryRepository {
	return &orderItemQueryRepository{
		queryClient:   queryClient,
		commandClient: commandClient,
	}
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.GetOrderItemsByOrderRow, error) {
	res, err := r.queryClient.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		return nil, order_item_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}

	var items []*db.GetOrderItemsByOrderRow
	for _, item := range res.Data {
		items = append(items, &db.GetOrderItemsByOrderRow{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})
	}

	return items, nil
}

func (r *orderItemQueryRepository) CalculateTotalPrice(ctx context.Context, order_id int) (*int32, error) {
	res, err := r.commandClient.CalculateTotalPrice(ctx, &pb.CalculateTotalPriceRequest{OrderId: int32(order_id)})
	if err != nil {
		return nil, order_item_errors.ErrCalculateTotalPrice.WithInternal(err)
	}

	total := int32(res.TotalPrice)
	return &total, nil
}
