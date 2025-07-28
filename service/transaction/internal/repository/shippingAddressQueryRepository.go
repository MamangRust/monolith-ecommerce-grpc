package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type shippingAddressQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.ShippingAddressMapping
}

func NewShippingAddressQueryRepository(db *db.Queries, mapping recordmapper.ShippingAddressMapping) *shippingAddressQueryRepository {
	return &shippingAddressQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *shippingAddressQueryRepository) FindByOrder(ctx context.Context, order_id int) (*record.ShippingAddressRecord, error) {
	res, err := r.db.GetShippingAddressByOrderID(ctx, int32(order_id))

	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}

	return r.mapping.ToShippingAddressRecord(res), nil

}
