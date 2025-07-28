package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ShippingAddressQueryRepository interface {
	FindAllShippingAddress(ctx context.Context, req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByOrder(ctx context.Context, shipping_id int) (*record.ShippingAddressRecord, error)
	FindById(ctx context.Context, shipping_id int) (*record.ShippingAddressRecord, error)
}

type ShippingAddressCommandRepository interface {
	CreateShippingAddress(ctx context.Context, request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	UpdateShippingAddress(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	TrashShippingAddress(ctx context.Context, category_id int) (*record.ShippingAddressRecord, error)
	RestoreShippingAddress(ctx context.Context, category_id int) (*record.ShippingAddressRecord, error)
	DeleteShippingAddressPermanently(ctx context.Context, category_id int) (bool, error)
	RestoreAllShippingAddress(ctx context.Context) (bool, error)
	DeleteAllPermanentShippingAddress(ctx context.Context) (bool, error)
}
