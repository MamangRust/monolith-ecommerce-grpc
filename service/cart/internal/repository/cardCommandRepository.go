package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type cartCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CartRecordMapping
}

func NewCartCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CartRecordMapping) *cartCommandRepository {
	return &cartCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cartCommandRepository) CreateCart(req *requests.CartCreateRecord) (*record.CartRecord, error) {
	res, err := r.db.CreateCart(r.ctx, db.CreateCartParams{
		UserID:    int32(req.UserID),
		ProductID: int32(req.ProductID),
		Name:      req.Name,
		Price:     int32(req.Price),
		Image:     req.ImageProduct,
		Quantity:  int32(req.Quantity),
		Weight:    int32(req.Weight),
	})

	if err != nil {
		return nil, cart_errors.ErrCreateCart
	}

	return r.mapping.ToCartRecord(res), nil
}

func (r *cartCommandRepository) DeletePermanent(req *requests.DeleteCartRequest) (bool, error) {
	err := r.db.DeleteCartByIdAndUserId(r.ctx, db.DeleteCartByIdAndUserIdParams{
		UserID: int32(req.UserID),
		CartID: int32(req.CartID),
	})

	if err != nil {
		return false, cart_errors.ErrDeleteCartPermanent
	}

	return true, nil
}

func (r *cartCommandRepository) DeleteAllPermanently(req *requests.DeleteAllCartRequest) (bool, error) {
	cartIDs := make([]int32, len(req.CartIds))

	for i, id := range req.CartIds {
		cartIDs[i] = int32(id)
	}

	err := r.db.DeleteAllCartByUserId(r.ctx, db.DeleteAllCartByUserIdParams{
		UserID:  int32(req.UserID),
		Column1: cartIDs,
	})

	if err != nil {
		return false, cart_errors.ErrDeleteAllCarts
	}

	return true, nil
}
