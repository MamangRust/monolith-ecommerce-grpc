package mencache

import (
	"context"
	"fmt"
)

type shippingAddressCommandCache struct {
	store *CacheStore
}

func NewShippingAddressCommandCache(store *CacheStore) *shippingAddressCommandCache {
	return &shippingAddressCommandCache{store: store}
}

func (shippingAddressCommandCache *shippingAddressCommandCache) DeleteShippingAddressCache(ctx context.Context, shipping_id int) {
	DeleteFromCache(ctx, shippingAddressCommandCache.store, fmt.Sprintf(shippingAddressByIdCacheKey, shipping_id))
}
