package merchantdetail_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type merchantDetailCommandCache struct {
	store *cache.CacheStore
}

func NewMerchantDetailCommandCache(store *cache.CacheStore) *merchantDetailCommandCache {
	return &merchantDetailCommandCache{store: store}
}

func (m *merchantDetailCommandCache) DeleteMerchantDetailCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantDetailByIdCacheKey, id)
	cache.DeleteFromCache(ctx, m.store, key)
}
