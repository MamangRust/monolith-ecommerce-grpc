package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"go.uber.org/zap"
)

type orderCommandCache struct {
	store *cache.CacheStore
}

func NewOrderCommandCache(store *cache.CacheStore) *orderCommandCache {
	return &orderCommandCache{store: store}
}

func (s *orderCommandCache) DeleteOrderCache(ctx context.Context, orderID int) {
	if _, err := s.store.InvalidateCache(ctx, "order:*"); err != nil {
		cache.DeleteFromCache(ctx, s.store, fmt.Sprintf(orderByIdCacheKey, orderID))
	}
}

func (s *orderCommandCache) InvalidateOrderCache(ctx context.Context) {
	if _, err := s.store.InvalidateCache(ctx, "order:*"); err != nil {
		s.store.Logger.Error("failed to invalidate order cache", zap.Error(err))
	}
}
