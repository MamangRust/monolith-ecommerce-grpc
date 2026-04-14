package cache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	OrderQueryCache           OrderQueryCache
	OrderCommandCache         OrderCommandCache
	OrderStatsCache           OrderStatsCache
	OrderStatsByMerchantCache OrderStatsByMerchantCache
}

func NewMencache(cacheStore *cache.CacheStore) *Mencache {
	return &Mencache{
		OrderQueryCache:           NewOrderQueryCache(cacheStore),
		OrderCommandCache:         NewOrderCommandCache(cacheStore),
		OrderStatsCache:           NewOrderStatsCache(cacheStore),
		OrderStatsByMerchantCache: NewOrderStatsByMerchantCache(cacheStore),
	}
}
