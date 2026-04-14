package review_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type ReviewMencache interface {
	QueryCache() ReviewQueryCache
	CommandCache() ReviewCommandCache
}

type reviewMencache struct {
	ReviewQueryCache
	ReviewCommandCache
}

func (m *reviewMencache) QueryCache() ReviewQueryCache {
	return m.ReviewQueryCache
}

func (m *reviewMencache) CommandCache() ReviewCommandCache {
	return m.ReviewCommandCache
}

func NewReviewMencache(cacheStore *cache.CacheStore) ReviewMencache {
	return &reviewMencache{
		ReviewQueryCache:   NewReviewQueryCache(cacheStore),
		ReviewCommandCache: NewReviewCommandCache(cacheStore),
	}
}
