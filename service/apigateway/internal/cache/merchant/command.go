package merchant_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
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
