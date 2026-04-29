package cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type Mencache struct {
	ReviewDetailQuery   ReviewDetailQueryCache
	ReviewDetailCommand ReviewDetailCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) *Mencache {
	return &Mencache{
		ReviewDetailQuery:   NewReviewDetailQueryCache(cacheStore),
		ReviewDetailCommand: NewReviewDetailCommandCache(cacheStore),
	}
}
