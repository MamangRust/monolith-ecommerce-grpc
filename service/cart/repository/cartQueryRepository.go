package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
)

type cartQueryRepository struct {
	db *db.Queries
}

func NewCartQueryRepository(db *db.Queries) CartQueryRepository {
	return &cartQueryRepository{
		db: db,
	}
}

func (r *cartQueryRepository) FindCarts(ctx context.Context, req *requests.FindAllCarts) ([]*db.GetCartsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCartsParams{
		UserID:  int32(req.UserID),
		Column2: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCarts(ctx, reqDb)

	if err != nil {
		return nil, cart_errors.ErrFindAllCarts.WithInternal(err)
	}

	return res, nil
}

