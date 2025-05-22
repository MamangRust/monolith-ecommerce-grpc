package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ShippingAddressQueryService interface {
	FindAll(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(shipping_id int) (*response.ShippingAddressResponse, *response.ErrorResponse)
	FindByOrder(order_id int) (*response.ShippingAddressResponse, *response.ErrorResponse)
}

type ShippingAddressCommandService interface {
	TrashShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse)
	RestoreShippingAddress(shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse)
	DeleteShippingAddressPermanently(categoryID int) (bool, *response.ErrorResponse)
	RestoreAllShippingAddress() (bool, *response.ErrorResponse)
	DeleteAllPermanentShippingAddress() (bool, *response.ErrorResponse)
}
