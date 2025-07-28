package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	ReviewDetailQueryCache   ReviewDetailQueryCache
	ReviewDetailCommandCache ReviewDetailCommandCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		ReviewDetailQueryCache:   NewReviewDetailQueryCache(cacheStore),
		ReviewDetailCommandCache: NewReviewDetailCommandCache(cacheStore),
	}
}
