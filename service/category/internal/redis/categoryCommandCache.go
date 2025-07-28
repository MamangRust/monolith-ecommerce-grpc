package mencache

import (
	"context"
	"fmt"
)

type categoryCommandCache struct {
	store *CacheStore
}

func NewCategoryCommandCache(store *CacheStore) *categoryCommandCache {
	return &categoryCommandCache{store: store}
}

func (c *categoryCommandCache) DeleteCachedCategoryCache(ctx context.Context, id int) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	DeleteFromCache(ctx, c.store, key)
}
