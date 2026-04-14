package repository

import (
	"context"

	"database/sql"
 
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
)


type categoryQueryRepository struct {
	db *db.Queries
}

func NewCategoryQueryRepository(db *db.Queries) *categoryQueryRepository {
	return &categoryQueryRepository{
		db: db,
	}
}

func (r *categoryQueryRepository) FindById(ctx context.Context, category_id int) (*db.GetCategoryByIDRow, error) {
	res, err := r.db.GetCategoryByID(ctx, int32(category_id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, product_errors.ErrProductNotFound.WithInternal(err)
		}
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return res, nil
}
