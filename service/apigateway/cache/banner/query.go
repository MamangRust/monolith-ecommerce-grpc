package banner_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

type bannerQueryCache struct {
	store *cache.CacheStore
}

func NewBannerQueryCache(store *cache.CacheStore) *bannerQueryCache {
	return &bannerQueryCache{store: store}
}

func (b *bannerQueryCache) GetCachedBanners(
	ctx context.Context,
	req *requests.FindAllBanner,
) (*response.ApiResponsePaginationBanner, bool) {

	key := fmt.Sprintf(
		bannerAllCacheKey,
		req.Page,
		req.PageSize,
		req.Search,
	)

	result, found := cache.GetFromCache[response.ApiResponsePaginationBanner](
		ctx,
		b.store,
		key,
	)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (b *bannerQueryCache) SetCachedBanners(
	ctx context.Context,
	req *requests.FindAllBanner,
	data *response.ApiResponsePaginationBanner,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(
		bannerAllCacheKey,
		req.Page,
		req.PageSize,
		req.Search,
	)

	cache.SetToCache(ctx, b.store, key, data, ttlDefault)
}

func (b *bannerQueryCache) GetCachedActiveBanners(
	ctx context.Context,
	req *requests.FindAllBanner,
) (*response.ApiResponsePaginationBannerDeleteAt, bool) {

	key := fmt.Sprintf(
		bannerActiveCacheKey,
		req.Page,
		req.PageSize,
		req.Search,
	)

	result, found := cache.GetFromCache[response.ApiResponsePaginationBannerDeleteAt](
		ctx,
		b.store,
		key,
	)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (b *bannerQueryCache) SetCachedActiveBanners(
	ctx context.Context,
	req *requests.FindAllBanner,
	data *response.ApiResponsePaginationBannerDeleteAt,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(
		bannerActiveCacheKey,
		req.Page,
		req.PageSize,
		req.Search,
	)

	cache.SetToCache(ctx, b.store, key, data, ttlDefault)
}

func (b *bannerQueryCache) GetCachedTrashedBanners(
	ctx context.Context,
	req *requests.FindAllBanner,
) (*response.ApiResponsePaginationBannerDeleteAt, bool) {

	key := fmt.Sprintf(
		bannerTrashedCacheKey,
		req.Page,
		req.PageSize,
		req.Search,
	)

	result, found := cache.GetFromCache[response.ApiResponsePaginationBannerDeleteAt](
		ctx,
		b.store,
		key,
	)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (b *bannerQueryCache) SetCachedTrashedBanners(
	ctx context.Context,
	req *requests.FindAllBanner,
	data *response.ApiResponsePaginationBannerDeleteAt,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(
		bannerTrashedCacheKey,
		req.Page,
		req.PageSize,
		req.Search,
	)

	cache.SetToCache(ctx, b.store, key, data, ttlDefault)
}

func (b *bannerQueryCache) GetCachedBanner(
	ctx context.Context,
	id int,
) (*response.ApiResponseBanner, bool) {

	key := fmt.Sprintf(bannerByIdCacheKey, id)

	result, found := cache.GetFromCache[response.ApiResponseBanner](
		ctx,
		b.store,
		key,
	)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (b *bannerQueryCache) SetCachedBanner(
	ctx context.Context,
	data *response.ApiResponseBanner,
) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(bannerByIdCacheKey, data.Data.ID)

	cache.SetToCache(ctx, b.store, key, data, ttlDefault)
}
