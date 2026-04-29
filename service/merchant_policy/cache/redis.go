package cache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	MerchantPoliciesCommandCache MerchantPoliciesCommandCache
	MerchantPoliciesQueryCache   MerchantPoliciesQueryCache
}

func NewMencache(cacheStore *cache.CacheStore) *Mencache {
	return &Mencache{
		MerchantPoliciesQueryCache:   NewMerchantPoliciesQueryCache(cacheStore),
		MerchantPoliciesCommandCache: NewMerchantPoliciesCommandCache(cacheStore),
	}
}
