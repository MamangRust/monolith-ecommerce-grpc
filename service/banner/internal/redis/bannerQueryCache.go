package mencache

import (
	"context"
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

func (b *bannerQueryCache) GetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponse, *int, bool) {
	key := fmt.Sprintf(bannerAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[bannerCacheResponse](ctx, b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner, data []*response.BannerResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	key := fmt.Sprintf(bannerAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponse{Data: data, Total: total}
	SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(bannerActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[bannerCacheResponseDeleteAt](ctx, b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner, data []*response.BannerResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.BannerResponseDeleteAt{}
	}

	key := fmt.Sprintf(bannerActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponseDeleteAt{Data: data, Total: total}
	SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(bannerTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[bannerCacheResponseDeleteAt](ctx, b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner, data []*response.BannerResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.BannerResponseDeleteAt{}
	}

	key := fmt.Sprintf(bannerTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponseDeleteAt{Data: data, Total: total}
	SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerCache(ctx context.Context, id int) (*response.BannerResponse, bool) {
	key := fmt.Sprintf(bannerByIdCacheKey, id)

	result, found := GetFromCache[*response.BannerResponse](ctx, b.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (b *bannerQueryCache) SetCachedBannerCache(ctx context.Context, data *response.BannerResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(bannerByIdCacheKey, data.ID)

	SetToCache(ctx, b.store, key, data, ttlDefault)
}
