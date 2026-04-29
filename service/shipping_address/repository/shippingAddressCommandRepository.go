package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
)

type shippingAddressCommandRepository struct {
	db *db.Queries
}

func NewShippingAddressCommandRepository(db *db.Queries) *shippingAddressCommandRepository {
	return &shippingAddressCommandRepository{
		db: db,
	}
}

func (r *shippingAddressCommandRepository) Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*db.CreateShippingAddressRow, error) {
	var orderID int32
	if request.OrderID != nil {
		orderID = int32(*request.OrderID)
	}
	req := db.CreateShippingAddressParams{
		OrderID: orderID,
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   float64(request.ShippingCost),
	}

	address, err := r.db.CreateShippingAddress(ctx, req)

	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}

	return address, nil
}

func (r *shippingAddressCommandRepository) Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*db.UpdateShippingAddressRow, error) {
	req := db.UpdateShippingAddressParams{
		ShippingAddressID: int32(*request.ShippingID),
		Alamat:            request.Alamat,
		Provinsi:          request.Provinsi,
		Kota:              request.Kota,
		Negara:            request.Negara,
		Courier:           request.Courier,
		ShippingMethod:    request.ShippingMethod,
		ShippingCost:      float64(request.ShippingCost),
	}

	res, err := r.db.UpdateShippingAddress(ctx, req)
	if err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress
	}

	return res, nil
}

func (r *shippingAddressCommandRepository) Trash(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	res, err := r.db.TrashShippingAddress(ctx, int32(shipping_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrTrashShippingAddress
	}

	return res, nil
}

func (r *shippingAddressCommandRepository) Restore(ctx context.Context, shipping_id int) (*db.ShippingAddress, error) {
	res, err := r.db.RestoreShippingAddress(ctx, int32(shipping_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrRestoreShippingAddress
	}

	return res, nil
}

func (r *shippingAddressCommandRepository) DeletePermanent(ctx context.Context, shipping_id int) (bool, error) {
	err := r.db.DeleteShippingAddressPermanently(ctx, int32(shipping_id))

	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}

	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, order_id int) (bool, error) {
	err := r.db.DeleteShippingAddressByOrderPermanent(ctx, int32(order_id))

	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}

	return true, nil
}

func (r *shippingAddressCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllShippingAddress(ctx)

	if err != nil {
		return false, shippingaddress_errors.ErrRestoreAllShippingAddresses
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentShippingAddress(ctx)

	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress
	}
	return true, nil
}
