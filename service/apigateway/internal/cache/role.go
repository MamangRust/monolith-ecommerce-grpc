package cache

import (
	"context"
	"fmt"
	"time"
	shared_cache "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

const (
	ttlDefault = 5 * time.Minute

	cacheRoleKey = "user_roles:%s"
)

type roleCache struct {
	store *shared_cache.CacheStore
}

func NewRoleCache(store *shared_cache.CacheStore) RoleCache {
	return &roleCache{store: store}
}

func (c *roleCache) GetRoleCache(ctx context.Context, userID string) ([]string, bool) {
	key := fmt.Sprintf(cacheRoleKey, userID)

	result, found := shared_cache.GetFromCache[[]string](ctx, c.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (c *roleCache) SetRoleCache(ctx context.Context, userID string, roles []string) {
	if userID == "" || len(roles) == 0 {
		return
	}

	key := fmt.Sprintf(cacheRoleKey, userID)

	shared_cache.SetToCache(ctx, c.store, key, &roles, ttlDefault)
}
