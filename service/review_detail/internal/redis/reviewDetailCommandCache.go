package mencache

import "fmt"

type reviewDetailCommandCache struct {
	store *CacheStore
}

func NewReviewDetailCommandCache(store *CacheStore) *reviewDetailCommandCache {
	return &reviewDetailCommandCache{store: store}
}

func (s *reviewDetailCommandCache) DeleteReviewDetailCache(review_id int) {
	DeleteFromCache(s.store, fmt.Sprintf(reviewByIdCacheKey, review_id))
}
