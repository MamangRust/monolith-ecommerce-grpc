package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CartQueryCache interface {
	GetCachedCartsCache(ctx context.Context, request *requests.FindAllCarts) ([]*db.GetCartsRow, *int, bool)
	SetCartsCache(ctx context.Context, request *requests.FindAllCarts, response []*db.GetCartsRow, total *int)
}
