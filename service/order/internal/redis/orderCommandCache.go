package mencache

import (
	"context"
	"fmt"
)

type orderCommandCache struct {
	store *CacheStore
}

func NewOrderCommandCache(store *CacheStore) *orderCommandCache {
	return &orderCommandCache{store: store}
}

func (s *orderCommandCache) DeleteOrderCache(ctx context.Context, order_id int) {
	DeleteFromCache(ctx, s.store, fmt.Sprintf(orderByIdCacheKey, order_id))
}
