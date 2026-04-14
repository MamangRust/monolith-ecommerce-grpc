package cart_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type cartMencache struct {
	CartQueryCache
}

type CartMencache interface {
	CartQueryCache
}

func NewCartMencache(cacheStore *cache.CacheStore) CartMencache {
	return &cartMencache{
		CartQueryCache: NewCartQueryCache(cacheStore),
	}
}
