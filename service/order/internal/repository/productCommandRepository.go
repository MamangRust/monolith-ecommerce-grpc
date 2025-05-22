package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type productCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.ProductRecordMapping
}

func NewProductCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ProductRecordMapping) *productCommandRepository {
	return &productCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *productCommandRepository) UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error) {
	res, err := r.db.UpdateProductCountStock(r.ctx, db.UpdateProductCountStockParams{
		ProductID:    int32(product_id),
		CountInStock: int32(stock),
	})

	if err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return r.mapping.ToProductRecord(res), nil
}
