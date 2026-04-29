package shippingaddress_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ShippingAddressQueryCache interface {
	GetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress) (*response.ApiResponsePaginationShippingAddress, bool)
	SetShippingAddressAllCache(ctx context.Context, req *requests.FindAllShippingAddress, data *response.ApiResponsePaginationShippingAddress)

	GetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress) (*response.ApiResponsePaginationShippingAddressDeleteAt, bool)
	SetShippingAddressActiveCache(ctx context.Context, req *requests.FindAllShippingAddress, data *response.ApiResponsePaginationShippingAddressDeleteAt)

	GetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress) (*response.ApiResponsePaginationShippingAddressDeleteAt, bool)
	SetShippingAddressTrashedCache(ctx context.Context, req *requests.FindAllShippingAddress, data *response.ApiResponsePaginationShippingAddressDeleteAt)

	GetCachedShippingAddressCache(ctx context.Context, shipping_id int) (*response.ApiResponseShippingAddress, bool)
	SetCachedShippingAddressCache(ctx context.Context, data *response.ApiResponseShippingAddress)

	GetCachedShippingAddressByOrderCache(ctx context.Context, order_id int) (*response.ApiResponseShippingAddress, bool)
	SetCachedShippingAddressByOrderCache(ctx context.Context, data *response.ApiResponseShippingAddress)
}

type ShippingAddressCommandCache interface {
	DeleteShippingAddressCache(ctx context.Context, shipping_id int)
}
