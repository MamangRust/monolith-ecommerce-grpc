package banner_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type bannerCommandCache struct {
	store *cache.CacheStore
}

func NewBannerCommandCache(store *cache.CacheStore) *bannerCommandCache {
	return &bannerCommandCache{store: store}
}

func (b *bannerCommandCache) DeleteBannerCache(ctx context.Context, id int) {
	key := fmt.Sprintf(bannerByIdCacheKey, id)

	cache.DeleteFromCache(ctx, b.store, key)
}
