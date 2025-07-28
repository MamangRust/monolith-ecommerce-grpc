package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	productAllCacheKey      = "product:all:page:%d:pageSize:%d:search:%s"
	productCategoryCacheKey = "product:category:%s:page:%d:pageSize:%d:search:%s"
	productMerchantCacheKey = "product:merchant:%d:page:%d:pageSize:%d:search:%s"

	productActiveCacheKey  = "product:active:page:%d:pageSize:%d:search:%s"
	productTrashedCacheKey = "product:trashed:page:%d:pageSize:%d:search:%s"
	productByIdCacheKey    = "product:id:%d"

	ttlDefault = 5 * time.Minute
)

type productCacheResponse struct {
	Data         []*response.ProductResponse `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type productCacheResponseDeleteAt struct {
	Data         []*response.ProductResponseDeleteAt `json:"data"`
	TotalRecords *int                                `json:"total_records"`
}

type productQueryCache struct {
	store *CacheStore
}

func NewProductQueryCache(store *CacheStore) *productQueryCache {
	return &productQueryCache{store: store}
}

func (p *productQueryCache) GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponse, *int, bool) {
	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[productCacheResponse](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ProductResponse{}
	}

	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*response.ProductResponse, *int, bool) {
	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[productCacheResponse](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant, data []*response.ProductResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ProductResponse{}
	}

	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*response.ProductResponse, *int, bool) {
	key := fmt.Sprintf(productCategoryCacheKey, req.CategoryName, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[productCacheResponse](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory, data []*response.ProductResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ProductResponse{}
	}

	key := fmt.Sprintf(productCategoryCacheKey, req.CategoryName, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[productCacheResponseDeleteAt](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ProductResponseDeleteAt{}
	}

	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[productCacheResponseDeleteAt](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ProductResponseDeleteAt{}
	}

	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProduct(ctx context.Context, productID int) (*response.ProductResponse, bool) {
	key := fmt.Sprintf(productByIdCacheKey, productID)

	result, found := GetFromCache[*response.ProductResponse](ctx, p.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (p *productQueryCache) SetCachedProduct(ctx context.Context, data *response.ProductResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(productByIdCacheKey, data.ID)
	SetToCache(ctx, p.store, key, data, ttlDefault)
}
