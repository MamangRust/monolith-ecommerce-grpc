package transaction_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type transactionCommandCache struct {
	store *cache.CacheStore
}

func NewTransactionCommandCache(store *cache.CacheStore) *transactionCommandCache {
	return &transactionCommandCache{store: store}
}

func (t *transactionCommandCache) DeleteTransactionCache(ctx context.Context, transactionID int) {
	key := fmt.Sprintf(transactionByIdCacheKey, transactionID)
	cache.DeleteFromCache(ctx, t.store, key)
}

func (t *transactionCommandCache) InvalidateTransactionCache(ctx context.Context) {
	t.store.InvalidateCache(ctx, "transaction:*")
}
