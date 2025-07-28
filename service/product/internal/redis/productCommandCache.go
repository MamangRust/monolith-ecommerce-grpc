package mencache

import (
	"context"
	"fmt"
)

type productCommandCache struct {
	store *CacheStore
}

func NewProductCommandCache(store *CacheStore) *productCommandCache {
	return &productCommandCache{store: store}
}

func (c *productCommandCache) DeleteCachedProduct(ctx context.Context, productID int) {
	DeleteFromCache(ctx, c.store, fmt.Sprintf(productByIdCacheKey, productID))
}
