package cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type sliderMencache struct {
	SliderQueryCache
	SliderCommandCache
}

type SliderMencache interface {
	SliderQueryCache
	SliderCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) SliderMencache {
	return sliderMencache{
		SliderQueryCache:   NewSliderQueryCache(cacheStore),
		SliderCommandCache: NewSliderCommandCache(cacheStore),
	}
}
