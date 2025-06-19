package mencache

import "fmt"

type merchantAwardCommandCache struct {
	store *CacheStore
}

func NewMerchantAwardCommandCache(store *CacheStore) *merchantAwardCommandCache {
	return &merchantAwardCommandCache{store: store}
}

func (m *merchantAwardCommandCache) DeleteMerchantAwardCache(id int) {
	key := fmt.Sprintf(merchantAwardByIdCacheKey, id)
	DeleteFromCache(m.store, key)
}
