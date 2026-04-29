package cart_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	cartAllCacheKey = "cart:all:page:%d:pageSize:%d:search:%s"
	ttlDefault      = 5 * time.Minute
)

type cartQueryCache struct {
	store *cache.CacheStore
}

func NewCartQueryCache(store *cache.CacheStore) *cartQueryCache {
	return &cartQueryCache{store: store}
}

func (c *cartQueryCache) GetCachedCarts(
	ctx context.Context,
	request *requests.FindAllCarts,
) (*response.ApiResponseCartPagination, bool) {

	key := fmt.Sprintf(
		cartAllCacheKey,
		request.Page,
		request.PageSize,
		request.Search,
	)

	result, found := cache.GetFromCache[response.ApiResponseCartPagination](
		ctx,
		c.store,
		key,
	)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (c *cartQueryCache) SetCachedCarts(
	ctx context.Context,
	request *requests.FindAllCarts,
	resp *response.ApiResponseCartPagination,
) {
	if resp == nil {
		return
	}

	key := fmt.Sprintf(
		cartAllCacheKey,
		request.Page,
		request.PageSize,
		request.Search,
	)

	cache.SetToCache(ctx, c.store, key, resp, ttlDefault)
}
