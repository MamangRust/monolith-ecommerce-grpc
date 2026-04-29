package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
)

type shippingAddressQueryRepository struct {
	db *db.Queries
}

func NewShippingAddressQueryRepository(db *db.Queries) *shippingAddressQueryRepository {
	return &shippingAddressQueryRepository{
		db: db,
	}
}

func (r *shippingAddressQueryRepository) FindAll(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddress(ctx, reqDb)

	if err != nil {
		return nil, shippingaddress_errors.ErrFindAllShippingAddress
	}

	return res, nil
}

func (r *shippingAddressQueryRepository) FindActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressActive(ctx, reqDb)

	if err != nil {
		return nil, shippingaddress_errors.ErrFindActiveShippingAddress
	}

	return res, nil
}

func (r *shippingAddressQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetShippingAddressTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetShippingAddressTrashed(ctx, reqDb)

	if err != nil {
		return nil, shippingaddress_errors.ErrFindTrashedShippingAddress
	}

	return res, nil
}

func (r *shippingAddressQueryRepository) FindByID(ctx context.Context, shipping_id int) (*db.GetShippingByIDRow, error) {
	res, err := r.db.GetShippingByID(ctx, int32(shipping_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByID
	}

	return res, nil
}

func (r *shippingAddressQueryRepository) FindByOrder(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, error) {
	res, err := r.db.GetShippingAddressByOrderID(ctx, int32(order_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}

	return res, nil
}
