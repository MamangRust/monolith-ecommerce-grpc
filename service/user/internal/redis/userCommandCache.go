package mencache

import "fmt"

type userCommandCache struct {
	store *CacheStore
}

func NewUserCommandCache(store *CacheStore) *userCommandCache {
	return &userCommandCache{store: store}
}

func (u *userCommandCache) DeleteUserCache(id int) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	DeleteFromCache(u.store, key)
}
