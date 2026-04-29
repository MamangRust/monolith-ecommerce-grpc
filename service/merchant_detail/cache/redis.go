package cache

import (
	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	MerchantDetailQueryCache   MerchantDetailQueryCache
	MerchantDetailCommandCache MerchantDetailCommandCache
}

func NewMencache(cacheStore *sharedcachehelpers.CacheStore) *Mencache {
	return &Mencache{
		MerchantDetailQueryCache:   NewMerchantDetailQueryCache(cacheStore),
		MerchantDetailCommandCache: NewMerchantDetailCommandCache(cacheStore),
	}
}
