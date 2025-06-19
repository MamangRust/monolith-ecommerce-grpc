package mencache

import "fmt"

type roleCommandCache struct {
	store *CacheStore
}

func NewRoleCommandCache(store *CacheStore) *roleCommandCache {
	return &roleCommandCache{store: store}
}

func (s *roleCommandCache) DeleteCachedRole(id int) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	DeleteFromCache(s.store, key)
}
