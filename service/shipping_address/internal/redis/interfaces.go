package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ShippingAddressQueryCache interface {
	GetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, bool)
	SetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponse, total *int)

	GetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, bool)
	SetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponseDeleteAt, total *int)

	GetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, bool)
	SetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponseDeleteAt, total *int)

	GetCachedShippingAddressCache(ctx context.Context, shipping_id int) (*response.ShippingAddressResponse, bool)
	SetCachedShippingAddressCache(ctx context.Context, data *response.ShippingAddressResponse)

	GetCachedShippingAddressByOrderCache(ctx context.Context, order_id int) (*response.ShippingAddressResponse, bool)
	SetCachedShippingAddressByOrderCache(ctx context.Context, data *response.ShippingAddressResponse)
}

type ShippingAddressCommandCache interface {
	DeleteShippingAddressCache(ctx context.Context, shipping_id int)
}
