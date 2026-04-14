package cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type ShippingAddressMencache interface {
	ShippingAddressQueryCache
	ShippingAddressCommandCache
}

type shippingAddressMencache struct {
	ShippingAddressQueryCache
	ShippingAddressCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) ShippingAddressMencache {
	return &shippingAddressMencache{
		ShippingAddressQueryCache:   NewShippingAddressQueryCache(cacheStore),
		ShippingAddressCommandCache: NewShippingAddressCommandCache(cacheStore),
	}
}
