package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ShippingAddressQueryService interface {
	FindAllShippingAddress(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressRow, *int, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressActiveRow, *int, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllShippingAddress,
	) ([]*db.GetShippingAddressTrashedRow, *int, error)

	FindByOrder(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingAddressByOrderIDRow, error)

	FindById(
		ctx context.Context,
		shipping_id int,
	) (*db.GetShippingByIDRow, error)
}

type ShippingAddressCommandService interface {
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
