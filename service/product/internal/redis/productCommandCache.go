package mencache

import "fmt"

type productCommandCache struct {
	store *CacheStore
}

func NewProductCommandCache(store *CacheStore) *productCommandCache {
	return &productCommandCache{store: store}
}

func (c *productCommandCache) DeleteCachedProduct(productID int) {
	DeleteFromCache(c.store, fmt.Sprintf(productByIdCacheKey, productID))
}
