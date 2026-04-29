package cache

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/user"
	category_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/category"
)


type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

type CacheApiGateway struct {
	RoleCache RoleCache
	UserCache user_cache.UserMencache
	CategoryCache category_cache.CategoryMencache
	CacheStore *cache.CacheStore
}

func NewCacheApiGateway(deps *Deps) *CacheApiGateway {
	cacheStore := cache.NewCacheStore(deps.Redis, deps.Logger, nil) // Metrics is nil for now
	return &CacheApiGateway{
		RoleCache: NewRoleCache(cacheStore),
		UserCache: user_cache.NewUserMencache(cacheStore),
		CategoryCache: category_cache.NewCategoryMencache(cacheStore),
		CacheStore: cacheStore,
	}
}
