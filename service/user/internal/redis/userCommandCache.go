package mencache

import (
	"context"
	"fmt"
)

type userCommandCache struct {
	store *CacheStore
}

func NewUserCommandCache(store *CacheStore) *userCommandCache {
	return &userCommandCache{store: store}
}

func (u *userCommandCache) DeleteUserCache(ctx context.Context, id int) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	DeleteFromCache(ctx, u.store, key)
}
