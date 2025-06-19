package mencache

import "fmt"

type categoryCommandCache struct {
	store *CacheStore
}

func NewCategoryCommandCache(store *CacheStore) *categoryCommandCache {
	return &categoryCommandCache{store: store}
}

func (c *categoryCommandCache) DeleteCachedCategoryCache(id int) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	DeleteFromCache(c.store, key)
}
