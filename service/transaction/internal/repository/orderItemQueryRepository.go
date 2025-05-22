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
	ctx     context.Context
	mapping recordmapper.OrderItemRecordMapping
}

func NewOrderItemQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.OrderItemRecordMapping) *orderItemQueryRepository {
	return &orderItemQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(order_id int) ([]*record.OrderItemRecord, error) {
	res, err := r.db.GetOrderItemsByOrder(r.ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	return r.mapping.ToOrderItemsRecord(res), nil
}
