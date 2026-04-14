package shippingaddress_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type ShippingAddressMencache interface {
	QueryCache() ShippingAddressQueryCache
	CommandCache() ShippingAddressCommandCache
}

type shippingAddressMencache struct {
	ShippingAddressQueryCache
	ShippingAddressCommandCache
}

func (m *shippingAddressMencache) QueryCache() ShippingAddressQueryCache {
	return m.ShippingAddressQueryCache
}

func (m *shippingAddressMencache) CommandCache() ShippingAddressCommandCache {
	return m.ShippingAddressCommandCache
}

func NewShippingAddressMencache(cacheStore *cache.CacheStore) ShippingAddressMencache {
	return &shippingAddressMencache{
		ShippingAddressQueryCache:   NewShippingAddressQueryCache(cacheStore),
		ShippingAddressCommandCache: NewShippingAddressCommandCache(cacheStore),
	}
}
