package cache

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

const (
	cartAllCacheKey = "cart:all:page:%d:pageSize:%d:search:%s"
	ttlDefault      = 5 * time.Minute
)

type cartCacheResponse struct {
	Data  []*db.GetCartsRow `json:"data"`
	Total *int              `json:"total_records"`
}

type cartQueryCache struct {
	store *cache.CacheStore
}

func NewCartQueryCache(store *cache.CacheStore) *cartQueryCache {
	return &cartQueryCache{store: store}
}

func (c *cartQueryCache) GetCachedCartsCache(ctx context.Context, request *requests.FindAllCarts) ([]*db.GetCartsRow, *int, bool) {
	key := fmt.Sprintf(cartAllCacheKey, request.Page, request.PageSize, request.Search)

	result, found := cache.GetFromCache[cartCacheResponse](ctx, c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *cartQueryCache) SetCartsCache(ctx context.Context, request *requests.FindAllCarts, response []*db.GetCartsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(cartAllCacheKey, request.Page, request.PageSize, request.Search)
	payload := &cartCacheResponse{Data: response, Total: total}
	cache.SetToCache(ctx, c.store, key, payload, ttlDefault)
}
