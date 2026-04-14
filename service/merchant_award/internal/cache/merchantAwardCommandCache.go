package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type merchantAwardCommandCache struct {
	store *cache.CacheStore
}

func NewMerchantAwardCommandCache(store *cache.CacheStore) *merchantAwardCommandCache {
	return &merchantAwardCommandCache{store: store}
}

func (m *merchantAwardCommandCache) DeleteMerchantAwardCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantAwardByIdCacheKey, id)
	cache.DeleteFromCache(ctx, m.store, key)
}
