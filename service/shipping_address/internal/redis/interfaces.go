package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ShippingAddressQueryCache interface {
	GetShippingAddressAllCache(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, bool)
	SetShippingAddressAllCache(req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponse, total *int)

	GetShippingAddressTrashedCache(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, bool)
	SetShippingAddressTrashedCache(req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponseDeleteAt, total *int)

	GetShippingAddressActiveCache(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, bool)
	SetShippingAddressActiveCache(req *requests.FindAllShippingAddress, res []*response.ShippingAddressResponseDeleteAt, total *int)

	GetCachedShippingAddressCache(shipping_id int) (*response.ShippingAddressResponse, bool)
	SetCachedShippingAddressCache(data *response.ShippingAddressResponse)

	GetCachedShippingAddressByOrderCache(order_id int) (*response.ShippingAddressResponse, bool)
	SetCachedShippingAddressByOrderCache(data *response.ShippingAddressResponse)
}

type ShippingAddressCommandCache interface {
	DeleteShippingAddressCache(shipping_id int)
}
