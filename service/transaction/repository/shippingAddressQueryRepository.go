package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type shippingAddressQueryRepository struct {
	client pb.ShippingQueryServiceClient
}

func NewShippingAddressQueryRepository(client pb.ShippingQueryServiceClient) *shippingAddressQueryRepository {
	return &shippingAddressQueryRepository{
		client: client,
	}
}

func (r *shippingAddressQueryRepository) FindByID(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, error) {
	res, err := r.client.FindByOrder(ctx, &pb.FindByIdShippingRequest{Id: int32(order_id)})
	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder.WithInternal(err)
	}

	return &db.GetShippingAddressByOrderIDRow{
		ShippingAddressID: res.Data.Id,
		OrderID:           res.Data.OrderId,
		Alamat:            res.Data.Alamat,
		Provinsi:          res.Data.Provinsi,
		Negara:            res.Data.Negara,
		Kota:              res.Data.Kota,
		ShippingMethod:    res.Data.ShippingMethod,
		ShippingCost:      float64(res.Data.ShippingCost),
	}, nil
}
