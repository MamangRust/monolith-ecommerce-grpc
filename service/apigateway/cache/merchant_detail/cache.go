package merchantdetail_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type MerchantDetailMencache interface {
	MerchantDetailQueryCache
	MerchantDetailCommandCache
}

type merchantDetailMencache struct {
	MerchantDetailQueryCache
	MerchantDetailCommandCache
}

func NewMerchantDetailMencache(cacheStore *cache.CacheStore) MerchantDetailMencache {
	return &merchantDetailMencache{
		MerchantDetailQueryCache:   NewMerchantDetailQueryCache(cacheStore),
		MerchantDetailCommandCache: NewMerchantDetailCommandCache(cacheStore),
	}
}
