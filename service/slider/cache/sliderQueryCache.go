package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	sliderAllCacheKey     = "slider:all:page:%d:pageSize:%d:search:%s"
	sliderActiveCacheKey  = "slider:active:page:%d:pageSize:%d:search:%s"
	sliderTrashedCacheKey = "slider:trashed:page:%d:pageSize:%d:search:%s"
	sliderIdKey           = "slider:id:%d"

	ttlDefault = 5 * time.Minute
)

type sliderCacheResponseDB struct {
	Data  []*db.GetSlidersRow `json:"data"`
	Total *int                `json:"total_records"`
}

type sliderActiveCacheResponseDB struct {
	Data  []*db.GetSlidersActiveRow `json:"data"`
	Total *int                      `json:"total_records"`
}

type sliderTrashedCacheResponseDB struct {
	Data  []*db.GetSlidersTrashedRow `json:"data"`
	Total *int                       `json:"total_records"`
}

type sliderQueryCache struct {
	store *cache.CacheStore
}

func NewSliderQueryCache(store *cache.CacheStore) *sliderQueryCache {
	return &sliderQueryCache{store: store}
}

func (s *sliderQueryCache) GetSliderAllCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersRow, *int, bool) {
	key := fmt.Sprintf(sliderAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[sliderCacheResponseDB](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderAllCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetSlidersRow{}
	}

	key := fmt.Sprintf(sliderAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *sliderQueryCache) GetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersActiveRow, *int, bool) {
	key := fmt.Sprintf(sliderActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[sliderActiveCacheResponseDB](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetSlidersActiveRow{}
	}

	key := fmt.Sprintf(sliderActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderActiveCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *sliderQueryCache) GetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersTrashedRow, *int, bool) {
	key := fmt.Sprintf(sliderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[sliderTrashedCacheResponseDB](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider, data []*db.GetSlidersTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetSlidersTrashedRow{}
	}

	key := fmt.Sprintf(sliderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderTrashedCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *sliderQueryCache) GetSliderCache(ctx context.Context, slider_id int) (*db.GetSliderByIDRow, bool) {
	key := fmt.Sprintf(sliderIdKey, slider_id)
	result, found := cache.GetFromCache[db.GetSliderByIDRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *sliderQueryCache) SetSliderCache(ctx context.Context, data *db.GetSliderByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(sliderIdKey, int(data.SliderID))
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
