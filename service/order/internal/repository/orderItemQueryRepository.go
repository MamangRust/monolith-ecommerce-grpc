package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
)

type orderItemQueryRepository struct {
	db *db.Queries
}

func NewOrderItemQueryRepository(db *db.Queries) OrderItemQueryRepository {
	return &orderItemQueryRepository{
		db: db,
	}
}

func (r *orderItemQueryRepository) CalculateTotalPrice(ctx context.Context, order_id int) (*int32, error) {
	res, err := r.db.CalculateTotalPrice(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice.WithInternal(err)
	}

	return &res, nil

}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.GetOrderItemsByOrderRow, error) {
	res, err := r.db.GetOrderItemsByOrder(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}

	return res, nil
}

