package cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type Mencache struct {
	TransactionQueryCache           TransactionQueryCache
	TransactionCommandCache         TransactionCommandCache
	TransactionStatsCache           TransactionStatsCache
	TransactionStatsByMerchantCache TransactionStatsByMerchantCache
}

type TransactionMencache interface {
	TransactionQueryCache
	TransactionCommandCache
	TransactionStatsCache
	TransactionStatsByMerchantCache
}

func NewMencache(cacheStore *cache.CacheStore) *Mencache {
	return &Mencache{
		TransactionQueryCache:           NewTransactionQueryCache(cacheStore),
		TransactionCommandCache:         NewTransactionCommandCache(cacheStore),
		TransactionStatsCache:           NewTransactionStatsCache(cacheStore),
		TransactionStatsByMerchantCache: NewTransactionStatsByMerchantCache(cacheStore),
	}
}
