package merchantpolicies_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type merchantPolicyCommandCache struct {
	store *cache.CacheStore
}

func NewMerchantPolicyCommandCache(store *cache.CacheStore) *merchantPolicyCommandCache {
	return &merchantPolicyCommandCache{store: store}
}

func (m *merchantPolicyCommandCache) DeleteMerchantPolicyCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, id)
	cache.DeleteFromCache(ctx, m.store, key)
}
