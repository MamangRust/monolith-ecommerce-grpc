package cache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type Mencache struct {
	MerchantQueryCache           MerchantQueryCache
	MerchantCommandCache         MerchantCommandCache
	MerchantDocumentQueryCache   MerchantDocumentQueryCache
	MerchantDocumentCommandCache MerchantDocumentCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) *Mencache {
	return &Mencache{
		MerchantQueryCache:           NewMerchantQueryCache(cacheStore),
		MerchantCommandCache:         NewMerchantCommandCache(cacheStore),
		MerchantDocumentQueryCache:   NewMerchantDocumentQueryCache(cacheStore),
		MerchantDocumentCommandCache: NewMerchantDocumentCommandCache(cacheStore),
	}
}
