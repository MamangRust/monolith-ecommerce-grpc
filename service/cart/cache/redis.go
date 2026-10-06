package cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type cartMencache struct {
	CartQueryCache
	CartCommandCache
}

type CartMencache interface {
	CartQueryCache
	CartCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) CartMencache {
	return &cartMencache{
		CartQueryCache:   NewCartQueryCache(cacheStore),
		CartCommandCache: NewCartCommandCache(cacheStore),
	}
}
