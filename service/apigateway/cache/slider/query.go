package slider_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	sliderAllCacheKey     = "slider:all:page:%d:pageSize:%d:search:%s"
	sliderActiveCacheKey  = "slider:active:page:%d:pageSize:%d:search:%s"
	sliderTrashedCacheKey = "slider:trashed:page:%d:pageSize:%d:search:%s"
	sliderIdKey           = "slider:id:%d"

	ttlDefault = 5 * time.Minute
)

type sliderQueryCache struct {
	store *cache.CacheStore
}

func NewSliderQueryCache(store *cache.CacheStore) *sliderQueryCache {
	return &sliderQueryCache{store: store}
}

func (s *sliderQueryCache) GetSliderAllCache(ctx context.Context, req *requests.FindAllSlider) (*response.ApiResponsePaginationSlider, bool) {
	key := fmt.Sprintf(sliderAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationSlider](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *sliderQueryCache) SetSliderAllCache(ctx context.Context, req *requests.FindAllSlider, data *response.ApiResponsePaginationSlider) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(sliderAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *sliderQueryCache) GetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider) (*response.ApiResponsePaginationSliderDeleteAt, bool) {
	key := fmt.Sprintf(sliderActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationSliderDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *sliderQueryCache) SetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider, data *response.ApiResponsePaginationSliderDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(sliderActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *sliderQueryCache) GetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider) (*response.ApiResponsePaginationSliderDeleteAt, bool) {
	key := fmt.Sprintf(sliderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationSliderDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *sliderQueryCache) SetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider, data *response.ApiResponsePaginationSliderDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(sliderTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *sliderQueryCache) GetCachedSliderCache(ctx context.Context, id int) (*response.ApiResponseSlider, bool) {
	key := fmt.Sprintf(sliderIdKey, id)

	result, found := cache.GetFromCache[response.ApiResponseSlider](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *sliderQueryCache) SetCachedSliderCache(ctx context.Context, data *response.ApiResponseSlider) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(sliderIdKey, data.Data.ID)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
