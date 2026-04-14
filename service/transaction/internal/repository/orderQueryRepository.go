package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
)

type orderQueryRepository struct {
	db *db.Queries
}

func NewOrderQueryRepository(db *db.Queries) OrderQueryRepository {
	return &orderQueryRepository{
		db: db,
	}
}

func (r *orderQueryRepository) FindById(ctx context.Context, id int) (*db.GetOrderByIDRow, error) {
	res, err := r.db.GetOrderByID(ctx, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, order_errors.ErrOrderNotFound.WithInternal(err)
		}
		return nil, order_errors.ErrFindById.WithInternal(err)
	}

	return res, nil
}

