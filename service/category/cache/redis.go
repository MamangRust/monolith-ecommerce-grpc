package cache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	CategoryQueryCache           CategoryQueryCache
	CategoryCommandCache         CategoryCommandCache
	CategoryStatsCache           CategoryStatsCache
	CategoryStatsByIdCache       CategoryStatsByIdCache
	CategoryStatsByMerchantCache CategoryStatsByMerchantCache
}

func NewMencache(cacheStore *cache.CacheStore) *Mencache {
	return &Mencache{
		CategoryQueryCache:           NewCategoryQueryCache(cacheStore),
		CategoryCommandCache:         NewCategoryCommandCache(cacheStore),
		CategoryStatsCache:           NewCategoryStatsCache(cacheStore),
		CategoryStatsByIdCache:       NewCategoryStatsByIdCache(cacheStore),
		CategoryStatsByMerchantCache: NewCategoryStatsByMerchantCache(cacheStore),
	}
}
