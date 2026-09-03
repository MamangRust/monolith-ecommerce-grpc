package cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type transactionCommandCache struct {
	store *cache.CacheStore
}

func NewTransactionCommandCache(store *cache.CacheStore) *transactionCommandCache {
	return &transactionCommandCache{store: store}
}

func (t *transactionCommandCache) DeleteTransactionCache(ctx context.Context, transactionID int) {
	if _, err := t.store.InvalidateCache(ctx, "transaction:*"); err != nil {
		cache.DeleteFromCache(ctx, t.store, fmt.Sprintf(transactionByIdCacheKey, transactionID))
	}
}

func (t *transactionCommandCache) InvalidateTransactionCache(ctx context.Context) {
	t.store.InvalidateCache(ctx, "transaction:*")
}
