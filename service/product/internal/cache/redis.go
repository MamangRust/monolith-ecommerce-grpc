package cache

import (
	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	ProductQuery   ProductQueryCache
	ProductCommand ProductCommandCache
}

func NewMencache(cacheStore *sharedcachehelpers.CacheStore) *Mencache {
	return &Mencache{
		ProductQuery:   NewProductQueryCache(cacheStore),
		ProductCommand: NewProductCommandCache(cacheStore),
	}
}
