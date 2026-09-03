package cache

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"go.uber.org/zap"
)

// cartSvcAllCacheKey is namespaced separately from the API gateway cart cache
// ("cart:all:...") because both layers share the same Redis but store different
// payload schemas on these keys.
const (
	cartSvcAllCacheKey = "cart:svc:all:user:%d:page:%d:pageSize:%d:search:%s"
	ttlDefault         = 5 * time.Minute
)

type cartCacheResponse struct {
	Data  []*db.GetCartsRow `json:"data"`
	Total *int              `json:"total_records"`
}

type cartQueryCache struct {
	store *cache.CacheStore
}

func NewCartQueryCache(store *cache.CacheStore) *cartQueryCache {
	return &cartQueryCache{store: store}
}

func (c *cartQueryCache) GetCachedCartsCache(ctx context.Context, request *requests.FindAllCarts) ([]*db.GetCartsRow, *int, bool) {
	key := fmt.Sprintf(cartSvcAllCacheKey, request.UserID, request.Page, request.PageSize, request.Search)

	result, found := cache.GetFromCache[cartCacheResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *cartQueryCache) SetCartsCache(ctx context.Context, request *requests.FindAllCarts, response []*db.GetCartsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(cartSvcAllCacheKey, request.UserID, request.Page, request.PageSize, request.Search)
	payload := &cartCacheResponse{Data: response, Total: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}

func (c *cartQueryCache) DeleteCartsCache(ctx context.Context, userID int) {
	pattern := fmt.Sprintf("cart:svc:all:user:%d:*", userID)
	if _, err := c.store.InvalidateCache(ctx, pattern); err != nil {
		c.store.Logger.Error("Failed to invalidate cart service cache", zap.Error(err))
	}
}
