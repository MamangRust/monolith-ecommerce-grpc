package product_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type ProductMencache interface {
	ProductQueryCache
	ProductCommandCache
}

type productMencache struct {
	ProductQueryCache
	ProductCommandCache
}

func NewProductMencache(cacheStore *cache.CacheStore) ProductMencache {
	return &productMencache{
		ProductQueryCache:   NewProductQueryCache(cacheStore),
		ProductCommandCache: NewProductCommandCache(cacheStore),
	}
}
