package shippingaddress_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type shippingAddressCommandCache struct {
	store *cache.CacheStore
}

func NewShippingAddressCommandCache(store *cache.CacheStore) *shippingAddressCommandCache {
	return &shippingAddressCommandCache{store: store}
}

func (shippingAddressCommandCache *shippingAddressCommandCache) DeleteShippingAddressCache(ctx context.Context, shipping_id int) {
	cache.DeleteFromCache(ctx, shippingAddressCommandCache.store, fmt.Sprintf(shippingAddressByIdCacheKey, shipping_id))
}
