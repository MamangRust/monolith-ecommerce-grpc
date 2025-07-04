package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	sliderAllCacheKey     = "slider:all:page:%d:pageSize:%d:search:%s"
	sliderActiveCacheKey  = "slider:active:page:%d:pageSize:%d:search:%s"
	sliderTrashedCacheKey = "slider:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type sliderCacheResposne struct {
	Data  []*response.SliderResponse `json:"data"`
	Total *int                       `json:"total_records"`
}

type sliderCacheResposneDeleteAt struct {
	Data  []*response.SliderResponseDeleteAt `json:"data"`
	Total *int                               `json:"total_records"`
}

type sliderQueryCache struct {
	store *CacheStore
}

func NewSliderQueryCache(store *CacheStore) *sliderQueryCache {
	return &sliderQueryCache{store: store}
}

func (s *sliderQueryCache) GetSliderAllCache(req *requests.FindAllSlider) ([]*response.SliderResponse, *int, bool) {
	key := fmt.Sprintf(sliderAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[sliderCacheResposne](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderAllCache(req *requests.FindAllSlider, data []*response.SliderResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.SliderResponse{}
	}

	key := fmt.Sprintf(sliderAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderCacheResposne{Data: data, Total: total}
	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *sliderQueryCache) GetSliderActiveCache(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(sliderActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[sliderCacheResposneDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderActiveCache(req *requests.FindAllSlider, data []*response.SliderResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.SliderResponseDeleteAt{}
	}

	key := fmt.Sprintf(sliderActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderCacheResposneDeleteAt{Data: data, Total: total}
	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *sliderQueryCache) GetSliderTrashedCache(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(sliderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[sliderCacheResposneDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (s *sliderQueryCache) SetSliderTrashedCache(req *requests.FindAllSlider, data []*response.SliderResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.SliderResponseDeleteAt{}
	}

	key := fmt.Sprintf(sliderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &sliderCacheResposneDeleteAt{Data: data, Total: total}
	SetToCache(s.store, key, payload, ttlDefault)
}
