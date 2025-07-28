package mencache

import (
	"context"
	"fmt"
)

type merchantAwardCommandCache struct {
	store *CacheStore
}

func NewMerchantAwardCommandCache(store *CacheStore) *merchantAwardCommandCache {
	return &merchantAwardCommandCache{store: store}
}

func (m *merchantAwardCommandCache) DeleteMerchantAwardCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantAwardByIdCacheKey, id)
	DeleteFromCache(ctx, m.store, key)
}
