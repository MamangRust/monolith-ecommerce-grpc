package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
)

type productQueryRepository struct {
	db *db.Queries
}

func NewProductQueryRepository(db *db.Queries) ProductQueryRepository {
	return &productQueryRepository{
		db: db,
	}
}

func (r *productQueryRepository) FindById(ctx context.Context, id int) (*db.GetProductByIDRow, error) {
	res, err := r.db.GetProductByID(ctx, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, product_errors.ErrProductNotFound.WithInternal(err)
		}
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return res, nil
}
