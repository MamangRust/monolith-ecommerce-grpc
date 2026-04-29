package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CartQueryRepository interface {
	FindCarts(
		ctx context.Context,
		req *requests.FindAllCarts,
	) ([]*db.GetCartsRow, error)
}

type CartCommandRepository interface {
	CreateCart(
		ctx context.Context,
		req *requests.CartCreateRecord,
	) (*db.CreateCartRow, error)

	DeletePermanent(
		ctx context.Context,
		req *requests.DeleteCartRequest,
	) (bool, error)

	DeleteAllPermanently(
		ctx context.Context,
		req *requests.DeleteAllCartRequest,
	) (bool, error)
}

type ProductQueryRepository interface {
	FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
}

type UserQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
}
