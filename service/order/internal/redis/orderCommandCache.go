package mencache

import "fmt"

type orderCommandCache struct {
	store *CacheStore
}

func NewOrderCommandCache(store *CacheStore) *orderCommandCache {
	return &orderCommandCache{store: store}
}

func (s *orderCommandCache) DeleteOrderCache(order_id int) {
	DeleteFromCache(s.store, fmt.Sprintf(orderByIdCacheKey, order_id))
}
