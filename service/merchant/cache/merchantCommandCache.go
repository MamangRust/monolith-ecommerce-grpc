package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"go.uber.org/zap"
)

type merchantCommandCache struct {
	store *cache.CacheStore
}

func NewMerchantCommandCache(store *cache.CacheStore) *merchantCommandCache {
	return &merchantCommandCache{store: store}

}

func (s *merchantCommandCache) DeleteCachedMerchant(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	cache.DeleteFromCache(ctx, s.store, key)
}

func (s *merchantCommandCache) InvalidateMerchantCache(ctx context.Context) {
	if s == nil || s.store == nil {
		return
	}
	if _, err := s.store.InvalidateCache(ctx, "merchant:*"); err != nil {
		s.store.Logger.Error("failed to invalidate merchant cache", zap.Error(err))
	}
}
