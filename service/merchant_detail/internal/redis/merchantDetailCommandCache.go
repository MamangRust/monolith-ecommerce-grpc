package mencache

import "fmt"

type merchantDetailCommandCache struct {
	store *CacheStore
}

func NewMerchantDetailCommandCache(store *CacheStore) *merchantDetailCommandCache {
	return &merchantDetailCommandCache{store: store}
}

func (m *merchantDetailCommandCache) DeleteMerchantDetailCache(id int) {
	key := fmt.Sprintf(merchantDetailByIdCacheKey, id)
	DeleteFromCache(m.store, key)
}
