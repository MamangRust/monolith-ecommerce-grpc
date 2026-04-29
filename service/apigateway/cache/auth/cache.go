package auth_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type authMencache struct {
	IdentityCache
	LoginCache
}

type AuthMencache interface {
	IdentityCache
	LoginCache
}

func NewMencache(cacheStore *cache.CacheStore) AuthMencache {
	return &authMencache{
		IdentityCache: NewIdentityCache(cacheStore),
		LoginCache:    NewLoginCache(cacheStore),
	}
}
