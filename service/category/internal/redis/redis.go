package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	CategoryQueryCache           CategoryQueryCache
	CategoryCommandCache         CategoryCommandCache
	CategoryStatsCache           CategoryStatsCache
	CategoryStatsByIdCache       CategoryStatsByIdCache
	CategoryStatsByMerchantCache CategoryStatsByMerchantCache
}

type Deps struct {
	Ctx    context.Context
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Ctx, deps.Redis, deps.Logger)

	return &Mencache{
		CategoryQueryCache:           NewCategoryQueryCache(cacheStore),
		CategoryCommandCache:         NewCategoryCommandCache(cacheStore),
		CategoryStatsCache:           NewCategoryStatsCache(cacheStore),
		CategoryStatsByIdCache:       NewCategoryStatsByIdCache(cacheStore),
		CategoryStatsByMerchantCache: NewCategoryStatsByMerchantCache(cacheStore),
	}
}
