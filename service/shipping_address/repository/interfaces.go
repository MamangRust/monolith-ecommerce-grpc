package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ShippingAddressQueryRepository interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressRow, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressActiveRow, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressTrashedRow, error)

	FindByOrder(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingAddressByOrderIDRow, error)

	FindByID(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingByIDRow, error)
}

type ShippingAddressCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateShippingAddressRequest,
	) (*db.CreateShippingAddressRow, error)

	Update(
		ctx context.Context,
		request *requests.UpdateShippingAddressRequest,
	) (*db.UpdateShippingAddressRow, error)

	Trash(
		ctx context.Context,
		shipping_id int,
	) (*db.ShippingAddress, error)

	Restore(
		ctx context.Context,
		shipping_id int,
	) (*db.ShippingAddress, error)

	DeletePermanent(
		ctx context.Context,
		shipping_id int,
	) (bool, error)

	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
