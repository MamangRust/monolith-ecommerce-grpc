package merchantpolicies_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type merchantPoliciesMencache struct {
	MerchantPolicyCommandCache
	MerchantPolicyQueryCache
}

type MerchantPoliciesMencache interface {
	MerchantPolicyCommandCache
	MerchantPolicyQueryCache
}

func NewMerchantPoliciesMencache(cacheStore *cache.CacheStore) MerchantPoliciesMencache {
	return &merchantPoliciesMencache{
		MerchantPolicyQueryCache:   NewMerchantPolicyQueryCache(cacheStore),
		MerchantPolicyCommandCache: NewMerchantPolicyCommandCache(cacheStore),
	}
}
