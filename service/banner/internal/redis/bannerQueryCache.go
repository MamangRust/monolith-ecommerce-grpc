package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	bannerAllCacheKey     = "banner:all:page:%d:pageSize:%d:search:%s"
	bannerByIdCacheKey    = "banner:id:%d"
	bannerActiveCacheKey  = "banner:active:page:%d:pageSize:%d:search:%s"
	bannerTrashedCacheKey = "banner:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type bannerCacheResponse struct {
	Data  []*response.BannerResponse `json:"data"`
	Total *int                       `json:"total"`
}

type bannerCacheResponseDeleteAt struct {
	Data  []*response.BannerResponseDeleteAt `json:"data"`
	Total *int                               `json:"total"`
}

type bannerQueryCache struct {
	store *CacheStore
}

func NewBannerQueryCache(store *CacheStore) *bannerQueryCache {
	return &bannerQueryCache{store: store}
}

func (b *bannerQueryCache) GetCachedBannersCache(req *requests.FindAllBanner) ([]*response.BannerResponse, *int, bool) {
	key := fmt.Sprintf(bannerAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[bannerCacheResponse](b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannersCache(req *requests.FindAllBanner, data []*response.BannerResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	key := fmt.Sprintf(bannerAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponse{Data: data, Total: total}
	SetToCache(b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerActiveCache(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(bannerActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[bannerCacheResponseDeleteAt](b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannerActiveCache(req *requests.FindAllBanner, data []*response.BannerResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.BannerResponseDeleteAt{}
	}

	key := fmt.Sprintf(bannerActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponseDeleteAt{Data: data, Total: total}
	SetToCache(b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerTrashedCache(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(bannerTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[bannerCacheResponseDeleteAt](b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannerTrashedCache(req *requests.FindAllBanner, data []*response.BannerResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.BannerResponseDeleteAt{}
	}

	key := fmt.Sprintf(bannerTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponseDeleteAt{Data: data, Total: total}
	SetToCache(b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerCache(id int) (*response.BannerResponse, bool) {
	key := fmt.Sprintf(bannerByIdCacheKey, id)

	result, found := GetFromCache[*response.BannerResponse](b.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (b *bannerQueryCache) SetCachedBannerCache(data *response.BannerResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(bannerByIdCacheKey, data.ID)

	SetToCache(b.store, key, data, ttlDefault)
}
