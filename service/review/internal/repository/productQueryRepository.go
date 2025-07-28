package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type productQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.ProductRecordMapping
}

func NewProductQueryRepository(db *db.Queries, mapping recordmapper.ProductRecordMapping) *productQueryRepository {
	return &productQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *productQueryRepository) FindById(ctx context.Context, user_id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByID(ctx, int32(user_id))

	if err != nil {
		return nil, product_errors.ErrFindById
	}

	return r.mapping.ToProductRecord(res), nil
}
