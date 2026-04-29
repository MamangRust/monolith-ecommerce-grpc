package cache

import (
	"context"

	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type orderItemCommandCache struct {
	cacheStore *sharedcachehelpers.CacheStore
}

func NewOrderItemCommandCache(cacheStore *sharedcachehelpers.CacheStore) *orderItemCommandCache {
	return &orderItemCommandCache{
		cacheStore: cacheStore,
	}
}

func (c *orderItemCommandCache) InvalidateOrderItemCache(ctx context.Context) error {
	_, err := c.cacheStore.InvalidateCache(ctx, "order_item:*")
	return err
}
