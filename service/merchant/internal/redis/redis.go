package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	MerchantQueryCache           MerchantQueryCache
	MerchantCommandCache         MerchantCommandCache
	MerchantDocumentQueryCache   MerchantDocumentQueryCache
	MerchantDocumentCommandCache MerchantDocumentCommandCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		MerchantQueryCache:           NewMerchantQueryCache(cacheStore),
		MerchantCommandCache:         NewMerchantCommandCache(cacheStore),
		MerchantDocumentQueryCache:   NewMerchantDocumentQueryCache(cacheStore),
		MerchantDocumentCommandCache: NewMerchantDocumentCommandCache(cacheStore),
	}
}
