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
	mapping recordmapper.OrderRecordMapping
}

func NewOrderQueryRepository(db *db.Queries, mapping recordmapper.OrderRecordMapping) *orderQueryRepository {
	return &orderQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *orderQueryRepository) FindById(ctx context.Context, id int) (*record.OrderRecord, error) {
	res, err := r.db.GetOrderByID(ctx, int32(id))

	if err != nil {
		return nil, order_errors.ErrFindById
	}

	return r.mapping.ToOrderRecord(res), nil
}
