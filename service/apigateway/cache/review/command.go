package review_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type reviewCommandCache struct {
	store *cache.CacheStore
}

func NewReviewCommandCache(store *cache.CacheStore) *reviewCommandCache {
	return &reviewCommandCache{store: store}
}

func (s *reviewCommandCache) DeleteReviewCache(ctx context.Context, review_id int) {
	key := fmt.Sprintf(reviewByIdCacheKey, review_id)

	cache.DeleteFromCache(ctx, s.store, key)
}
