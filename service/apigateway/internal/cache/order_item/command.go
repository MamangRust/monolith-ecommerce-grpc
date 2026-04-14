package orderitem_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type orderItemCommandCache struct {
	store *cache.CacheStore
}

func NewOrderItemCommandCache(store *cache.CacheStore) *orderItemCommandCache {
	return &orderItemCommandCache{store: store}
}

func (c *orderItemCommandCache) DeleteCachedOrderItemByOrderId(ctx context.Context, orderId int) {
	key := fmt.Sprintf(orderItemByOrderCacheKey, orderId)

	cache.DeleteFromCache(ctx, c.store, key)
}
