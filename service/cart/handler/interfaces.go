package handler

import (
	pbcart "github.com/MamangRust/monolith-ecommerce-pb/cart"
)

type CartQueryHandler interface {
	pbcart.CartQueryServiceServer
}

type CartCommandHandler interface {
	pbcart.CartCommandServiceServer
}
