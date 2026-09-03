package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CartQueryCache interface {
	GetCachedCartsCache(ctx context.Context, request *requests.FindAllCarts) ([]*db.GetCartsRow, *int, bool)
	SetCartsCache(ctx context.Context, request *requests.FindAllCarts, response []*db.GetCartsRow, total *int)

	// DeleteCartsCache invalidates every cached cart listing for a user so
	// mutations (create/delete/delete-all) never serve stale data.
	DeleteCartsCache(ctx context.Context, userID int)
}
