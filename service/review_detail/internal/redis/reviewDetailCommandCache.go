package mencache

import (
	"context"
	"fmt"
)

type reviewDetailCommandCache struct {
	store *CacheStore
}

func NewReviewDetailCommandCache(store *CacheStore) *reviewDetailCommandCache {
	return &reviewDetailCommandCache{store: store}
}

func (s *reviewDetailCommandCache) DeleteReviewDetailCache(ctx context.Context, review_id int) {
	DeleteFromCache(ctx, s.store, fmt.Sprintf(reviewByIdCacheKey, review_id))
}
