package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ShippingAddressQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*response.ShippingAddressResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, shipping_id int) (*response.ShippingAddressResponse, *response.ErrorResponse)
	FindByOrder(ctx context.Context, order_id int) (*response.ShippingAddressResponse, *response.ErrorResponse)
}

type ShippingAddressCommandService interface {
	TrashShippingAddress(ctx context.Context, shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse)
	RestoreShippingAddress(ctx context.Context, shipping_id int) (*response.ShippingAddressResponseDeleteAt, *response.ErrorResponse)
	DeleteShippingAddressPermanently(ctx context.Context, categoryID int) (bool, *response.ErrorResponse)
	RestoreAllShippingAddress(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllPermanentShippingAddress(ctx context.Context) (bool, *response.ErrorResponse)
}
