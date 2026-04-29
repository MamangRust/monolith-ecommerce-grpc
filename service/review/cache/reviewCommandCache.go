package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type reviewCommandCache struct {
	store *cache.CacheStore
}

func NewReviewCommandCache(store *cache.CacheStore) ReviewCommandCache {
	return &reviewCommandCache{
		store: store,
	}
}

func (c *reviewCommandCache) DeleteReviewCache(ctx context.Context, reviewID int) {
	cache.DeleteFromCache(ctx, c.store, fmt.Sprintf("review:%d", reviewID))
	cache.DeleteFromCache(ctx, c.store, "reviews:all")
	cache.DeleteFromCache(ctx, c.store, "reviews:active")
	cache.DeleteFromCache(ctx, c.store, "reviews:trashed")
}
