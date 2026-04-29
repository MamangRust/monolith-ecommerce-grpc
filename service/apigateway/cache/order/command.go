package order_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type orderCommandCache struct {
	store *cache.CacheStore
}

func NewOrderCommandCache(store *cache.CacheStore) *orderCommandCache {
	return &orderCommandCache{store: store}
}

func (s *orderCommandCache) DeleteOrderCache(ctx context.Context, order_id int) {
	cache.DeleteFromCache(ctx, s.store, fmt.Sprintf(orderByIdCacheKey, order_id))
}
