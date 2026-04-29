package cache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type reviewDetailCommandCache struct {
	store *cache.CacheStore
}

func NewReviewDetailCommandCache(store *cache.CacheStore) ReviewDetailCommandCache {
	return &reviewDetailCommandCache{
		store: store,
	}
}

func (c *reviewDetailCommandCache) DeleteReviewDetailCache(ctx context.Context, reviewDetailID int) {
	cache.DeleteFromCache(ctx, c.store, "review_detail:all")
	cache.DeleteFromCache(ctx, c.store, "review_detail:active")
	cache.DeleteFromCache(ctx, c.store, "review_detail:trashed")
}
