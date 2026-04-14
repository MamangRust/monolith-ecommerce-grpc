package merchantawards_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type MerchantAwardMencache interface {
	MerchantAwardQueryCache
	MerchantAwardCommandCache
}

type merchantAwardMencache struct {
	MerchantAwardQueryCache
	MerchantAwardCommandCache
}

func NewMerchantAward(cacheStore *cache.CacheStore) MerchantAwardMencache {
	return &merchantAwardMencache{
		MerchantAwardQueryCache:   NewMerchantAwardQueryCache(cacheStore),
		MerchantAwardCommandCache: NewMerchantAwardCommandCache(cacheStore),
	}
}
