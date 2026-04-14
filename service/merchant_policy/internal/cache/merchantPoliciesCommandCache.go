package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type merchantPoliciesCommandCache struct {
	store *cache.CacheStore
}

func NewMerchantPoliciesCommandCache(store *cache.CacheStore) MerchantPoliciesCommandCache {
	return &merchantPoliciesCommandCache{store: store}
}

func (m *merchantPoliciesCommandCache) DeleteMerchantPolicyCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, id)
	cache.DeleteFromCache(ctx, m.store, key)
}
