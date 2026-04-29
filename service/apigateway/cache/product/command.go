package product_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type productCommandCache struct {
	store *cache.CacheStore
}

func NewProductCommandCache(store *cache.CacheStore) *productCommandCache {
	return &productCommandCache{store: store}
}

func (c *productCommandCache) DeleteCachedProduct(ctx context.Context, productID int) {
	cache.DeleteFromCache(ctx, c.store, fmt.Sprintf(productByIdCacheKey, productID))
}
