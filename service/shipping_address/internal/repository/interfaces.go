package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ShippingAddressQueryRepository interface {
	FindAllShippingAddress(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressTrashedRow, error)

	FindByOrder(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingAddressByOrderIDRow, error)

	FindById(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingByIDRow, error)
}

type ShippingAddressCommandRepository interface {
	CreateShippingAddress(
		ctx context.Context,
		request *requests.CreateShippingAddressRequest,
	) (*db.CreateShippingAddressRow, error)

	UpdateShippingAddress(
		ctx context.Context,
		request *requests.UpdateShippingAddressRequest,
	) (*db.UpdateShippingAddressRow, error)

	TrashShippingAddress(
		ctx context.Context,
		shipping_id int,
	) (*db.ShippingAddress, error)

	RestoreShippingAddress(
		ctx context.Context,
		shipping_id int,
	) (*db.ShippingAddress, error)

	DeleteShippingAddressPermanently(
		ctx context.Context,
		shipping_id int,
	) (bool, error)

	RestoreAllShippingAddress(ctx context.Context) (bool, error)
	DeleteAllPermanentShippingAddress(ctx context.Context) (bool, error)
}
