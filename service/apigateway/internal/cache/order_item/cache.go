package orderitem_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type OrderItemMencache interface {
	OrderItemQueryCache
	OrderItemCommandCache
}

type orderItemMencache struct {
	OrderItemQueryCache
	OrderItemCommandCache
}

func NewOrderItemMencache(cacheStore *cache.CacheStore) OrderItemMencache {
	return &orderItemMencache{
		OrderItemQueryCache:   NewOrderItemQueryCache(cacheStore),
		OrderItemCommandCache: NewOrderItemCommandCache(cacheStore),
	}
}
