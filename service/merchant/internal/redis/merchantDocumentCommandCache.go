package mencache

import (
	"context"
	"fmt"
)

type merchantDocumentCommandCache struct {
	store *CacheStore
}

func NewMerchantDocumentCommandCache(store *CacheStore) *merchantDocumentCommandCache {
	return &merchantDocumentCommandCache{store: store}
}

func (s *merchantDocumentCommandCache) DeleteCachedMerchantDocuments(ctx context.Context, id int) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)
	DeleteFromCache(ctx, s.store, key)
}
