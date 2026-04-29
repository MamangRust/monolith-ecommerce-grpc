package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CartQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCarts) ([]*db.GetCartsRow, *int, error)
}

type CartCommandService interface {
	Create(ctx context.Context, req *requests.CreateCartRequest) (*db.Cart, error)
	DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error)
	DeleteAll(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error)
}
