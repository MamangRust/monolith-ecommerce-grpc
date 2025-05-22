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
	ctx     context.Context
	mapping recordmapper.ProductRecordMapping
}

func NewProductQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ProductRecordMapping) *productQueryRepository {
	return &productQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *productQueryRepository) FindById(product_id int) (*record.ProductRecord, error) {
	res, err := r.db.GetProductByID(r.ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrFindById
	}

	return r.mapping.ToProductRecord(res), nil
}
