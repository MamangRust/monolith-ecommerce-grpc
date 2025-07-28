package mencache

import (
	"context"
	"fmt"
)

type bannerCommandCache struct {
	store *CacheStore
}

func NewBannerCommandCache(store *CacheStore) *bannerCommandCache {
	return &bannerCommandCache{store: store}
}

func (b *bannerCommandCache) DeleteBannerCache(ctx context.Context, id int) {
	key := fmt.Sprintf(bannerByIdCacheKey, id)

	DeleteFromCache(ctx, b.store, key)
}
