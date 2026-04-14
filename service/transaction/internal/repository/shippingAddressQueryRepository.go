package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
)

type shippingAddressQueryRepository struct {
	db *db.Queries
}

func NewShippingAddressQueryRepository(db *db.Queries) ShippingAddressQueryRepository {
	return &shippingAddressQueryRepository{
		db: db,
	}
}

func (r *shippingAddressQueryRepository) FindByOrder(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, error) {
	res, err := r.db.GetShippingAddressByOrderID(ctx, int32(order_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shippingaddress_errors.ErrShippingAddressNotFound.WithInternal(err)
		}
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder.WithInternal(err)
	}

	return res, nil
}

