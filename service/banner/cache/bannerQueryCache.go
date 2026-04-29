package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type bannerCacheResponseDB struct {
	Data  []*db.GetBannersRow `json:"data"`
	Total *int                `json:"total"`
}

type bannerActiveCacheResponseDB struct {
	Data  []*db.GetBannersActiveRow `json:"data"`
	Total *int                      `json:"total"`
}

type bannerTrashedCacheResponseDB struct {
	Data  []*db.GetBannersTrashedRow `json:"data"`
	Total *int                       `json:"total"`
}

type bannerQueryCache struct {
	store *cache.CacheStore
}

func NewBannerQueryCache(store *cache.CacheStore) *bannerQueryCache {
	return &bannerQueryCache{store: store}
}

func (b *bannerQueryCache) GetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, *int, bool) {
	key := fmt.Sprintf(bannerAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[bannerCacheResponseDB](ctx, b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner, data []*db.GetBannersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(bannerAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, *int, bool) {
	key := fmt.Sprintf(bannerActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[bannerActiveCacheResponseDB](ctx, b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner, data []*db.GetBannersActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetBannersActiveRow{}
	}

	key := fmt.Sprintf(bannerActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerActiveCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, *int, bool) {
	key := fmt.Sprintf(bannerTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[bannerTrashedCacheResponseDB](ctx, b.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (b *bannerQueryCache) SetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner, data []*db.GetBannersTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetBannersTrashedRow{}
	}

	key := fmt.Sprintf(bannerTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &bannerTrashedCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, b.store, key, payload, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBannerCache(ctx context.Context, id int) (*db.GetBannerRow, bool) {
	key := fmt.Sprintf(bannerByIdCacheKey, id)

	result, found := cache.GetFromCache[db.GetBannerRow](ctx, b.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (b *bannerQueryCache) SetCachedBannerCache(ctx context.Context, data *db.GetBannerRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(bannerByIdCacheKey, data.BannerID)

	cache.SetToCache(ctx, b.store, key, data, ttlDefault)
}
