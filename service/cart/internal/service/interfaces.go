package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CartQueryService interface {
	FindAll(req *requests.FindAllCarts) ([]*response.CartResponse, *int, *response.ErrorResponse)
}

type CartCommandService interface {
	CreateCart(req *requests.CreateCartRequest) (*response.CartResponse, *response.ErrorResponse)
	DeletePermanent(req *requests.DeleteCartRequest) (bool, *response.ErrorResponse)
	DeleteAllPermanently(req *requests.DeleteAllCartRequest) (bool, *response.ErrorResponse)
}
