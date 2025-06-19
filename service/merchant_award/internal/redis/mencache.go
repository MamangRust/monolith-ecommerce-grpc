package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	MerchantAwardQueryCache   MerchantAwardQueryCache
	MerchantAwardCommandCache MerchanrAwardCommandCache
}

type Deps struct {
	Ctx    context.Context
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Ctx, deps.Redis, deps.Logger)

	return &Mencache{
		MerchantAwardQueryCache:   NewMerchantAwardQueryCache(cacheStore),
		MerchantAwardCommandCache: NewMerchantAwardCommandCache(cacheStore),
	}
}
