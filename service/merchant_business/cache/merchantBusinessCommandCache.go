package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type merchantBusinessCommandCache struct {
	store *cache.CacheStore
}

func NewMerchantBusinessCommandCache(store *cache.CacheStore) *merchantBusinessCommandCache {
	return &merchantBusinessCommandCache{store: store}
}

func (m *merchantBusinessCommandCache) DeleteMerchantBusinessCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantBusinessByIdCacheKey, id)
	cache.DeleteFromCache(ctx, m.store, key)
}
