package mencache

import "fmt"

type bannerCommandCache struct {
	store *CacheStore
}

func NewBannerCommandCache(store *CacheStore) *bannerCommandCache {
	return &bannerCommandCache{store: store}
}

func (b *bannerCommandCache) DeleteBannerCache(id int) {
	key := fmt.Sprintf(bannerByIdCacheKey, id)

	DeleteFromCache(b.store, key)
}
