package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CartQueryCache interface {
	GetCachedCartsCache(request *requests.FindAllCarts) ([]*response.CartResponse, *int, bool)
	SetCartsCache(request *requests.FindAllCarts, response []*response.CartResponse, total *int)
}
