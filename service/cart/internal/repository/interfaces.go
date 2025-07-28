package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CartQueryRepository interface {
	FindCarts(ctx context.Context, req *requests.FindAllCarts) ([]*record.CartRecord, *int, error)
}

type CartCommandRepository interface {
	CreateCart(ctx context.Context, req *requests.CartCreateRecord) (*record.CartRecord, error)
	DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error)
	DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error)
}

type ProductQueryRepository interface {
	FindById(ctx context.Context, productID int) (*record.ProductRecord, error)
}

type UserQueryRepository interface {
	FindById(ctx context.Context, userID int) (*record.UserRecord, error)
}
