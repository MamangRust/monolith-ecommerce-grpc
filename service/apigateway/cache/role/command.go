package role_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type roleCommandCache struct {
	store *cache.CacheStore
}

func NewRoleCommandCache(store *cache.CacheStore) *roleCommandCache {
	return &roleCommandCache{store: store}
}

func (s *roleCommandCache) DeleteCachedRole(ctx context.Context, id int) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	cache.DeleteFromCache(ctx, s.store, key)
}
