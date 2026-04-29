package category_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type CategoryMencache interface {
	CategoryQueryCache
	CategoryCommandCache
	CategoryStatsCache
	CategoryStatsByIdCache
	CategoryStatsByMerchantCache
}

type categoryMencache struct {
	CategoryQueryCache
	CategoryCommandCache
	CategoryStatsCache
	CategoryStatsByIdCache
	CategoryStatsByMerchantCache
}

func NewCategoryMencache(cacheStore *cache.CacheStore) CategoryMencache {
	return &categoryMencache{
		CategoryQueryCache:           NewCategoryQueryCache(cacheStore),
		CategoryCommandCache:         NewCategoryCommandCache(cacheStore),
		CategoryStatsCache:           NewCategoryStatsCache(cacheStore),
		CategoryStatsByIdCache:       NewCategoryStatsByIdCache(cacheStore),
		CategoryStatsByMerchantCache: NewCategoryStatsByMerchantCache(cacheStore),
	}
}
