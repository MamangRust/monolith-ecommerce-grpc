package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CartQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCarts) ([]*response.CartResponse, *int, *response.ErrorResponse)
}

type CartCommandService interface {
	CreateCart(ctx context.Context, req *requests.CreateCartRequest) (*response.CartResponse, *response.ErrorResponse)
	DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, *response.ErrorResponse)
	DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, *response.ErrorResponse)
}
