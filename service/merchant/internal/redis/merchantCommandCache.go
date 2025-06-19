package mencache

import "fmt"

type merchantCommandCache struct {
	store *CacheStore
}

func NewMerchantCommandCache(store *CacheStore) *merchantCommandCache {
	return &merchantCommandCache{store: store}

}

func (s *merchantCommandCache) DeleteCachedMerchant(id int) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	DeleteFromCache(s.store, key)
}
