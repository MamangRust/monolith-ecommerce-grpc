package mencache

import "fmt"

type transactionCommandCache struct {
	store *CacheStore
}

func NewTransactionCommandCache(store *CacheStore) *transactionCommandCache {
	return &transactionCommandCache{store: store}
}

func (t *transactionCommandCache) DeleteTransactionCache(transactionID int) {
	key := fmt.Sprintf(transactionByIdCacheKey, transactionID)

	DeleteFromCache(t.store, key)
}
