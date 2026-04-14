package cache

import (
	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	MerchantAwardQueryCache   MerchantAwardQueryCache
	MerchantAwardCommandCache MerchantAwardCommandCache
}

func NewMencache(cacheStore *sharedcachehelpers.CacheStore) *Mencache {
	return &Mencache{
		MerchantAwardQueryCache:   NewMerchantAwardQueryCache(cacheStore),
		MerchantAwardCommandCache: NewMerchantAwardCommandCache(cacheStore),
	}
}
