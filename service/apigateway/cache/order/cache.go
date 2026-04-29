package order_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type orderMencache struct {
	OrderQueryCache
	OrderCommandCache
	OrderStatsCache
	OrderStatsByMerchantCache
}

type OrderMencache interface {
	OrderQueryCache
	OrderCommandCache
	OrderStatsCache
	OrderStatsByMerchantCache
}

func OrderNewMencache(cacheStore *cache.CacheStore) OrderMencache {
	return &orderMencache{
		OrderQueryCache:           NewOrderQueryCache(cacheStore),
		OrderCommandCache:         NewOrderCommandCache(cacheStore),
		OrderStatsCache:           NewOrderStatsCache(cacheStore),
		OrderStatsByMerchantCache: NewOrderStatsByMerchantCache(cacheStore),
	}
}
