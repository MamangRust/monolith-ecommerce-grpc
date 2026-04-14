package cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type cartMencache struct {
	CartQueryCache
}

type CartMencache interface {
	CartQueryCache
}

func NewMencache(cacheStore *cache.CacheStore) CartMencache {
	return &cartMencache{
		CartQueryCache: NewCartQueryCache(cacheStore),
	}
}
