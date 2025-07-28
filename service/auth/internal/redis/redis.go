package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	IdentityCache      IdentityCache
	LoginCache         LoginCache
	PasswordResetCache PasswordResetCache
	RegisterCache      RegisterCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		IdentityCache:      NewidentityCache(cacheStore),
		LoginCache:         NewLoginCache(cacheStore),
		PasswordResetCache: NewPasswordResetCache(cacheStore),
		RegisterCache:      NewRegisterCache(cacheStore),
	}
}
