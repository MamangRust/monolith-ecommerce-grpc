package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type userCommandCache struct {
	store *cache.CacheStore
}

func NewUserCommandCache(store *cache.CacheStore) UserCommandCache {
	return &userCommandCache{store: store}
}

func (u *userCommandCache) DeleteUserCache(ctx context.Context, id int) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	cache.DeleteFromCache(ctx, u.store, key)
}
