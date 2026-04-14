package user_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type UserMencache interface {
	UserQueryCache
	UserCommandCache
}

type usermencache struct {
	UserQueryCache
	UserCommandCache
}

func NewUserMencache(cacheStore *cache.CacheStore) UserMencache {
	return &usermencache{
		UserQueryCache:   NewUserQueryCache(cacheStore),
		UserCommandCache: NewUserCommandCache(cacheStore),
	}
}
