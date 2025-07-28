package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type orderItemQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.OrderItemRecordMapping
}

func NewOrderItemQueryRepository(db *db.Queries, mapping recordmapper.OrderItemRecordMapping) *orderItemQueryRepository {
	return &orderItemQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *orderItemQueryRepository) CalculateTotalPrice(ctx context.Context, order_id int) (*int32, error) {
	res, err := r.db.CalculateTotalPrice(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
	}

	return &res, nil

}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*record.OrderItemRecord, error) {
	res, err := r.db.GetOrderItemsByOrder(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	return r.mapping.ToOrderItemsRecord(res), nil
}
