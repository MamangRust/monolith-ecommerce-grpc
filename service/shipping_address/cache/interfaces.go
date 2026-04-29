package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ShippingAddressQueryCache interface {
	GetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressRow, *int, bool)
	SetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*db.GetShippingAddressRow, total *int)

	GetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressTrashedRow, *int, bool)
	SetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*db.GetShippingAddressTrashedRow, total *int)

	GetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressActiveRow, *int, bool)
	SetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*db.GetShippingAddressActiveRow, total *int)

	GetCachedShippingAddressCache(ctx context.Context, shipping_id int) (*db.GetShippingByIDRow, bool)
	SetCachedShippingAddressCache(ctx context.Context, data *db.GetShippingByIDRow)

	GetCachedShippingAddressByOrderCache(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, bool)
	SetCachedShippingAddressByOrderCache(ctx context.Context, data *db.GetShippingAddressByOrderIDRow)
}

type ShippingAddressCommandCache interface {
	DeleteShippingAddressCache(ctx context.Context, shipping_id int)
	InvalidateShippingAddressCache(ctx context.Context)
}
