package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ShippingAddressQueryRepository interface {
	FindAllShippingAddress(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByActive(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByTrashed(req *requests.FindAllShippingAddress) ([]*record.ShippingAddressRecord, *int, error)
	FindByOrder(shipping_id int) (*record.ShippingAddressRecord, error)
	FindById(shipping_id int) (*record.ShippingAddressRecord, error)
}

type ShippingAddressCommandRepository interface {
	CreateShippingAddress(request *requests.CreateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	UpdateShippingAddress(request *requests.UpdateShippingAddressRequest) (*record.ShippingAddressRecord, error)
	TrashShippingAddress(category_id int) (*record.ShippingAddressRecord, error)
	RestoreShippingAddress(category_id int) (*record.ShippingAddressRecord, error)
	DeleteShippingAddressPermanently(category_id int) (bool, error)
	RestoreAllShippingAddress() (bool, error)
	DeleteAllPermanentShippingAddress() (bool, error)
}
