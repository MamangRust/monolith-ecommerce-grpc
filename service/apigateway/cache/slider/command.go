package slider_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"fmt"
)

type sliderCommandCache struct {
	store *cache.CacheStore
}

func NewSliderCommandCache(store *cache.CacheStore) *sliderCommandCache {
	return &sliderCommandCache{
		store: store,
	}
}

func (s *sliderCommandCache) DeleteSliderCache(ctx context.Context, slider_id int) {
	key := fmt.Sprintf(sliderIdKey, slider_id)

	cache.DeleteFromCache(ctx, s.store, key)
}
