package mencache

import (
	"context"
	"fmt"
)

type roleCommandCache struct {
	store *CacheStore
}

func NewRoleCommandCache(store *CacheStore) *roleCommandCache {
	return &roleCommandCache{store: store}
}

func (s *roleCommandCache) DeleteCachedRole(ctx context.Context, id int) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	DeleteFromCache(ctx, s.store, key)
}
