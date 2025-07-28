package mencache

import (
	"context"
	"fmt"
)

type transactionCommandCache struct {
	store *CacheStore
}

func NewTransactionCommandCache(store *CacheStore) *transactionCommandCache {
	return &transactionCommandCache{store: store}
}

func (t *transactionCommandCache) DeleteTransactionCache(ctx context.Context, transactionID int) {
	key := fmt.Sprintf(transactionByIdCacheKey, transactionID)

	DeleteFromCache(ctx, t.store, key)
}
