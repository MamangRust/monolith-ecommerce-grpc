package mencache

import (
	"context"
	"fmt"
)

type merchantCommandCache struct {
	store *CacheStore
}

func NewMerchantCommandCache(store *CacheStore) *merchantCommandCache {
	return &merchantCommandCache{store: store}

}

func (s *merchantCommandCache) DeleteCachedMerchant(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	DeleteFromCache(ctx, s.store, key)
}
