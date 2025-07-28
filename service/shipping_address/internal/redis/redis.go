package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	ShippingAddressQueryCache   ShippingAddressQueryCache
	ShippingAddressCommandCache ShippingAddressCommandCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		ShippingAddressQueryCache:   NewShippingAddressQueryCache(cacheStore),
		ShippingAddressCommandCache: NewShippingAddressCommandCache(cacheStore),
	}
}
