package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
)

type productQueryRepository struct {
	db *db.Queries
}

func NewProductQueryRepository(db *db.Queries) *productQueryRepository {
	return &productQueryRepository{
		db: db,
	}
}

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error) {
	res, err := r.db.GetProductByID(ctx, int32(product_id))

	if err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return res, nil
}
