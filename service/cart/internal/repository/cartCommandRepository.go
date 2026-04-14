package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
)

type cartCommandRepository struct {
	db *db.Queries
}

func NewCartCommandRepository(db *db.Queries) CartCommandRepository {
	return &cartCommandRepository{
		db: db,
	}
}

func (r *cartCommandRepository) CreateCart(ctx context.Context, req *requests.CartCreateRecord) (*db.CreateCartRow, error) {
	res, err := r.db.CreateCart(ctx, db.CreateCartParams{
		UserID:    int32(req.UserID),
		ProductID: int32(req.ProductID),
		Name:      req.Name,
		Price:     int32(req.Price),
		Image:     req.ImageProduct,
		Quantity:  int32(req.Quantity),
		Weight:    int32(req.Weight),
	})

	if err != nil {
		return nil, cart_errors.ErrCreateCart.WithInternal(err)
	}

	return res, nil
}

func (r *cartCommandRepository) DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error) {
	err := r.db.DeleteCartByIdAndUserId(ctx, db.DeleteCartByIdAndUserIdParams{
		CartID: int32(req.CartID),
		UserID: int32(req.UserID),
	})

	if err != nil {
		return false, cart_errors.ErrDeleteCartPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *cartCommandRepository) DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error) {
	cartIDs := make([]int32, len(req.CartIds))

	for i, id := range req.CartIds {
		cartIDs[i] = int32(id)
	}

	err := r.db.DeleteAllCartByUserId(ctx, db.DeleteAllCartByUserIdParams{
		Column1: cartIDs,
		UserID:  int32(req.UserID),
	})

	if err != nil {
		return false, cart_errors.ErrDeleteAllCarts.WithInternal(err)
	}

	return true, nil
}

