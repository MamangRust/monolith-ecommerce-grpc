package cart_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.uber.org/zap"
)

const (
	cartAllCacheKey = "cart:all:user:%d:page:%d:pageSize:%d:search:%s"
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
		request.UserID,
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
		request.UserID,
		request.Page,
		request.PageSize,
		request.Search,
	)

	cache.SetToCache(ctx, c.store, key, resp, ttlDefault)
}

func (c *cartQueryCache) DeleteCachedCarts(ctx context.Context, userID int) {
	pattern := fmt.Sprintf("cart:all:user:%d:*", userID)
	if _, err := c.store.InvalidateCache(ctx, pattern); err != nil {
		c.store.Logger.Error("Failed to invalidate cart cache", zap.Error(err))
	}
}
