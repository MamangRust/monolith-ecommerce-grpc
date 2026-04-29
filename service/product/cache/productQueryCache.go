package cache

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
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

type productCacheResponseDB struct {
	Data         []*db.GetProductsRow `json:"data"`
	TotalRecords *int                 `json:"total_records"`
}

type productMerchantCacheResponseDB struct {
	Data         []*db.GetProductsByMerchantRow `json:"data"`
	TotalRecords *int                           `json:"total_records"`
}

type productCategoryCacheResponseDB struct {
	Data         []*db.GetProductsByCategoryNameRow `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type productActiveCacheResponseDB struct {
	Data         []*db.GetProductsActiveRow `json:"data"`
	TotalRecords *int                       `json:"total_records"`
}

type productTrashedCacheResponseDB struct {
	Data         []*db.GetProductsTrashedRow `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type productQueryCache struct {
	store *cache.CacheStore
}

func NewProductQueryCache(store *cache.CacheStore) *productQueryCache {
	return &productQueryCache{store: store}
}

func (p *productQueryCache) GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, *int, bool) {
	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productCacheResponseDB](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data []*db.GetProductsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetProductsRow{}
	}

	key := fmt.Sprintf(productAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, *int, bool) {
	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productMerchantCacheResponseDB](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant, data []*db.GetProductsByMerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetProductsByMerchantRow{}
	}

	key := fmt.Sprintf(productMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &productMerchantCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, *int, bool) {
	key := fmt.Sprintf(productCategoryCacheKey, req.CategoryName, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productCategoryCacheResponseDB](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory, data []*db.GetProductsByCategoryNameRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetProductsByCategoryNameRow{}
	}

	key := fmt.Sprintf(productCategoryCacheKey, req.CategoryName, req.Page, req.PageSize, req.Search)
	payload := &productCategoryCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, *int, bool) {
	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productActiveCacheResponseDB](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data []*db.GetProductsActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetProductsActiveRow{}
	}

	key := fmt.Sprintf(productActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, *int, bool) {
	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[productTrashedCacheResponseDB](ctx, p.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (p *productQueryCache) SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data []*db.GetProductsTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetProductsTrashedRow{}
	}

	key := fmt.Sprintf(productTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &productTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, p.store, key, payload, ttlDefault)
}

func (p *productQueryCache) GetCachedProduct(ctx context.Context, productID int) (*db.GetProductByIDRow, bool) {
	key := fmt.Sprintf(productByIdCacheKey, productID)
	result, found := cache.GetFromCache[db.GetProductByIDRow](ctx, p.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (p *productQueryCache) SetCachedProduct(ctx context.Context, data *db.GetProductByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(productByIdCacheKey, data.ProductID)
	cache.SetToCache(ctx, p.store, key, data, ttlDefault)
}
