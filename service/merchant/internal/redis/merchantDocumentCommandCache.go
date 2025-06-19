package mencache

import "fmt"

type merchantDocumentCommandCache struct {
	store *CacheStore
}

func NewMerchantDocumentCommandCache(store *CacheStore) *merchantDocumentCommandCache {
	return &merchantDocumentCommandCache{store: store}
}

func (s *merchantDocumentCommandCache) DeleteCachedMerchantDocuments(id int) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)
	DeleteFromCache(s.store, key)
}
