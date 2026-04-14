package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
)

type orderItemQueryRepository struct {
	db *db.Queries
}

func NewOrderItemQueryRepository(db *db.Queries) OrderItemRepository {
	return &orderItemQueryRepository{
		db: db,
	}
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.GetOrderItemsByOrderRow, error) {
	res, err := r.db.GetOrderItemsByOrder(ctx, int32(order_id))

	if err != nil {
		return nil, order_errors.ErrOrderItemNotFound.WithInternal(err)
	}

	return res, nil
}

