package slider_cache

import "github.com/MamangRust/monolith-ecommerce-shared/cache"

type SliderMencache interface {
	QueryCache() SliderQueryCache
	CommandCache() SliderCommandCache
}

type sliderMencache struct {
	SliderQueryCache
	SliderCommandCache
}

func (m sliderMencache) QueryCache() SliderQueryCache {
	return m.SliderQueryCache
}

func (m sliderMencache) CommandCache() SliderCommandCache {
	return m.SliderCommandCache
}

func NewSliderMencache(cacheStore *cache.CacheStore) SliderMencache {
	return &sliderMencache{
		SliderQueryCache:   NewSliderQueryCache(cacheStore),
		SliderCommandCache: NewSliderCommandCache(cacheStore),
	}
}
