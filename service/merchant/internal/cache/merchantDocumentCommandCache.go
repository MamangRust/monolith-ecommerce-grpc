package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type merchantDocumentCommandCache struct {
	store *cache.CacheStore
}

func NewMerchantDocumentCommandCache(store *cache.CacheStore) *merchantDocumentCommandCache {
	return &merchantDocumentCommandCache{store: store}
}

func (s *merchantDocumentCommandCache) DeleteCachedMerchantDocuments(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)
	cache.DeleteFromCache(ctx, s.store, key)
}
