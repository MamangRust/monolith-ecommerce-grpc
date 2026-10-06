package cache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type cartCommandCache struct {
	store *cache.CacheStore
}

func NewCartCommandCache(store *cache.CacheStore) *cartCommandCache {
	return &cartCommandCache{store: store}
}

func (c *cartCommandCache) InvalidateCartsCache(ctx context.Context) {
	// Invalidate common lists/patterns
	c.store.InvalidateCache(ctx, "cart:*")
}
