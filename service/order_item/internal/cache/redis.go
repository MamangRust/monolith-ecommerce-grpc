package cache

import (
	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	OrderItemQueryCache   OrderItemQueryCache
	OrderItemCommandCache OrderItemCommandCache
}

func NewMencache(cacheStore *sharedcachehelpers.CacheStore) *Mencache {

	return &Mencache{
		OrderItemQueryCache:   NewOrderItemQueryCache(cacheStore),
		OrderItemCommandCache: NewOrderItemCommandCache(cacheStore),
	}
}
