package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

type CacheApiGateway struct {
	RoleCache RoleCache
}

func NewCacheApiGateway(deps *Deps) *CacheApiGateway {
	store := NewCacheStore(deps.Redis, deps.Logger)

	return &CacheApiGateway{
		RoleCache: NewRoleCache(store),
	}
}
