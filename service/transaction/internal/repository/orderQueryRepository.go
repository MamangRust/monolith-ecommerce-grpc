package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type orderQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.OrderRecordMapping
}

func NewOrderQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.OrderRecordMapping) *orderQueryRepository {
	return &orderQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *orderQueryRepository) FindById(id int) (*record.OrderRecord, error) {
	res, err := r.db.GetOrderByID(r.ctx, int32(id))

	if err != nil {
		return nil, order_errors.ErrFindById
	}

	return r.mapping.ToOrderRecord(res), nil
}
