package transaction_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type transactionMencache struct {
	TransactionQueryCache
	TransactionCommandCache
	TransactionStatsCache
	TransactionStatsByMerchantCache
}

type TransactionMencache interface {
	TransactionQueryCache
	TransactionCommandCache
	TransactionStatsCache
	TransactionStatsByMerchantCache
}

func NewTransactionMencache(cacheStore *cache.CacheStore) *transactionMencache {
	return &transactionMencache{
		TransactionQueryCache:           NewTransactionQueryCache(cacheStore),
		TransactionCommandCache:         NewTransactionCommandCache(cacheStore),
		TransactionStatsCache:           NewTransactionStatsCache(cacheStore),
		TransactionStatsByMerchantCache: NewTransactionStatsByMerchantCache(cacheStore),
	}
}
