package mencache

import "fmt"

type shippingAddressCommandCache struct {
	store *CacheStore
}

func NewShippingAddressCommandCache(store *CacheStore) *shippingAddressCommandCache {
	return &shippingAddressCommandCache{store: store}
}

func (shippingAddressCommandCache *shippingAddressCommandCache) DeleteShippingAddressCache(shipping_id int) {
	DeleteFromCache(shippingAddressCommandCache.store, fmt.Sprintf(shippingAddressByIdCacheKey, shipping_id))
}
