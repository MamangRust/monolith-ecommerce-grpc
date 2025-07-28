package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
	"golang.org/x/net/context"
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

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByID(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrFindById
	}

	return r.mapping.ToProductRecord(res), nil
}
