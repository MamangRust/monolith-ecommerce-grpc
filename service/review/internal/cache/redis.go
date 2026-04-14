package cache

import (
	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	ReviewQuery   ReviewQueryCache
	ReviewCommand ReviewCommandCache
}

func NewMencache(cacheStore *sharedcachehelpers.CacheStore) *Mencache {
	return &Mencache{
		ReviewQuery:   NewReviewQueryCache(cacheStore),
		ReviewCommand: NewReviewCommandCache(cacheStore),
	}
}
