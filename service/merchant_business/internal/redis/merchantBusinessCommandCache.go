package mencache

import (
	"context"
	"fmt"
)

type merchantBusinessCommandCache struct {
	store *CacheStore
}

func NewMerchantBusinessCommandCache(store *CacheStore) *merchantBusinessCommandCache {
	return &merchantBusinessCommandCache{store: store}
}

func (m *merchantBusinessCommandCache) DeleteMerchantBusinessCache(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantBusinessByIdCacheKey, id)
	DeleteFromCache(ctx, m.store, key)
}
