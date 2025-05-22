package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CartQueryRepository interface {
	FindCarts(req *requests.FindAllCarts) ([]*record.CartRecord, *int, error)
}

type CartCommandRepository interface {
	CreateCart(req *requests.CartCreateRecord) (*record.CartRecord, error)
	DeletePermanent(cart_id int) (bool, error)
	DeleteAllPermanently(req *requests.DeleteCartRequest) (bool, error)
}

type ProductQueryRepository interface {
	FindById(product_id int) (*record.ProductRecord, error)
}

type UserQueryRepository interface {
	FindById(user_id int) (*record.UserRecord, error)
}
