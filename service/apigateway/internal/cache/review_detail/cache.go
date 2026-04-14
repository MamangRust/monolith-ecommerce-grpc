package reviewdetail_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type ReviewDetailMencache interface {
	QueryCache() ReviewDetailQueryCache
	CommandCache() ReviewDetailCommandCache
}

type reviewDetaiMencache struct {
	ReviewDetailQueryCache
	ReviewDetailCommandCache
}

func (m *reviewDetaiMencache) QueryCache() ReviewDetailQueryCache {
	return m.ReviewDetailQueryCache
}

func (m *reviewDetaiMencache) CommandCache() ReviewDetailCommandCache {
	return m.ReviewDetailCommandCache
}

func NewReviewDetailMencache(cacheStore *cache.CacheStore) ReviewDetailMencache {
	return &reviewDetaiMencache{
		ReviewDetailQueryCache:   NewReviewDetailQueryCache(cacheStore),
		ReviewDetailCommandCache: NewReviewDetailCommandCache(cacheStore),
	}
}
