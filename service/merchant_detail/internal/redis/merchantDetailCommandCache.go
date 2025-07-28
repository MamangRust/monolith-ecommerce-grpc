package mencache

import (
	"context"
	"fmt"
)

type merchantDetailCommandCache struct {
	store *CacheStore
}

func NewMerchantDetailCommandCache(store *CacheStore) *merchantDetailCommandCache {
	return &merchantDetailCommandCache{store: store}
}

func (m *merchantDetailCommandCache) DeleteMerchantDetailCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantDetailByIdCacheKey, id)
	DeleteFromCache(ctx, m.store, key)
}
