package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type shippingAddressCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.ShippingAddressMapping
}

func NewShippingAddressCommandRepository(db *db.Queries, mapping recordmapper.ShippingAddressMapping) *shippingAddressCommandRepository {
	return &shippingAddressCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *shippingAddressCommandRepository) CreateShippingAddress(ctx context.Context, request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error) {
	req := db.CreateShippingAddressParams{
		OrderID:        int32(*request.OrderID),
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

	return r.mapping.ToShippingAddressRecord(address), nil
}

func (r *shippingAddressCommandRepository) UpdateShippingAddress(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error) {
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

	return r.mapping.ToShippingAddressRecord(res), nil
}
