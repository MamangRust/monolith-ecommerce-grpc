package category_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type categoryCommandCache struct {
	store *cache.CacheStore
}

func NewCategoryCommandCache(store *cache.CacheStore) *categoryCommandCache {
	return &categoryCommandCache{store: store}
}

func (c *categoryCommandCache) DeleteCachedCategoryCache(ctx context.Context, id int) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	cache.DeleteFromCache(ctx, c.store, key)
}
