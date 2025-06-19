package mencache

import "fmt"

type merchantBusinessCommandCache struct {
	store *CacheStore
}

func NewMerchantBusinessCommandCache(store *CacheStore) *merchantBusinessCommandCache {
	return &merchantBusinessCommandCache{store: store}
}

func (m *merchantBusinessCommandCache) DeleteMerchantBusinessCache(id int) {
	key := fmt.Sprintf(merchantBusinessByIdCacheKey, id)
	DeleteFromCache(m.store, key)
}
