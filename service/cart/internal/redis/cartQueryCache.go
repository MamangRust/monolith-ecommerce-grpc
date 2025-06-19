package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	cartAllCacheKey = "cart:all:page:%d:pageSize:%d:search:%s"
	ttlDefault      = 5 * time.Minute
)

type cartCacheResponse struct {
	Data  []*response.CartResponse `json:"data"`
	Total *int                     `json:"total_records"`
}

type cartQueryCache struct {
	store *CacheStore
}

func NewCartQueryCache(store *CacheStore) *cartQueryCache {
	return &cartQueryCache{store: store}
}

func (c *cartQueryCache) GetCachedCartsCache(request *requests.FindAllCarts) ([]*response.CartResponse, *int, bool) {
	key := fmt.Sprintf(cartAllCacheKey, request.Page, request.PageSize, request.Search)

	result, found := GetFromCache[cartCacheResponse](c.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (c *cartQueryCache) SetCartsCache(request *requests.FindAllCarts, response []*response.CartResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	key := fmt.Sprintf(cartAllCacheKey, request.Page, request.PageSize, request.Search)
	payload := &cartCacheResponse{Data: response, Total: total}
	SetToCache(c.store, key, payload, ttlDefault)
}
