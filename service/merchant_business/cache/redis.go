package cache

import (
	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	MerchantBusinessQueryCache   MerchantBusinessQueryCache
	MerchantBusinessCommandCache MerchantBusinessCommandCache
}

func NewMencache(cacheStore *sharedcachehelpers.CacheStore) *Mencache {
	return &Mencache{
		MerchantBusinessQueryCache:   NewMerchantBusinessQueryCache(cacheStore),
		MerchantBusinessCommandCache: NewMerchantBusinessCommandCache(cacheStore),
	}
}
