package mencache

import "fmt"

type merchantPolicyCommandCache struct {
	store *CacheStore
}

func NewMerchantPolicyCommandCache(store *CacheStore) *merchantPolicyCommandCache {
	return &merchantPolicyCommandCache{store: store}
}

func (m *merchantPolicyCommandCache) DeleteMerchantPolicyCache(id int) {
	key := fmt.Sprintf(merchantPolicyByIdCacheKey, id)
	DeleteFromCache(m.store, key)
}
