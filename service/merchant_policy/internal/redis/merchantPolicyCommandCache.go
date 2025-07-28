package mencache

import (
	"context"
	"fmt"
)

type merchantPolicyCommandCache struct {
	store *CacheStore
}

func NewMerchantPolicyCommandCache(store *CacheStore) *merchantPolicyCommandCache {
	return &merchantPolicyCommandCache{store: store}
}

func (m *merchantPolicyCommandCache) DeleteMerchantPolicyCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, id)
	DeleteFromCache(ctx, m.store, key)
}
