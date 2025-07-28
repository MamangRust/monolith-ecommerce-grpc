package mencache

import (
	"context"
	"fmt"
	"time"
)

const (
	ttlDefault = 5 * time.Minute

	cacheRoleKey = "user_roles:%s"
)

type roleCache struct {
	store *CacheStore
}

func NewRoleCache(store *CacheStore) RoleCache {
	return &roleCache{store: store}
}

func (c *roleCache) GetRoleCache(ctx context.Context, userID string) ([]string, bool) {
	key := fmt.Sprintf(cacheRoleKey, userID)

	result, found := GetFromCache[[]string](ctx, c.store, key)
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

	SetToCache(ctx, c.store, key, &roles, ttlDefault)
}
