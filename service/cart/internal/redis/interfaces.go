package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type CartQueryCache interface {
	GetCachedCartsCache(ctx context.Context, request *requests.FindAllCarts) ([]*response.CartResponse, *int, bool)
	SetCartsCache(ctx context.Context, request *requests.FindAllCarts, response []*response.CartResponse, total *int)
}
