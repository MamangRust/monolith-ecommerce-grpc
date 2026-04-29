package merchant_cache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type MerchantMencache interface {
	MerchantQueryCache
	MerchantCommandCache
}

type merchantMencache struct {
	MerchantQueryCache
	MerchantCommandCache
}

func NewMerchantMencache(cacheStore *cache.CacheStore) MerchantMencache {
	return &merchantMencache{
		MerchantQueryCache:   NewMerchantQueryCache(cacheStore),
		MerchantCommandCache: NewMerchantCommandCache(cacheStore),
	}
}
