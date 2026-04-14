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
	merchantDetailAllCacheKey     = "merchant_detail:all:page:%d:pageSize:%d:search:%s"
	merchantDetailByIdCacheKey    = "merchant_detail:id:%d"
	merchantDetailActiveCacheKey  = "merchant_detail:active:page:%d:pageSize:%d:search:%s"
	merchantDetailTrashedCacheKey = "merchant_detail:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantDetailCacheResponseDB struct {
	Data         []*db.GetMerchantDetailsRow `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type merchantDetailActiveCacheResponseDB struct {
	Data         []*db.GetMerchantDetailsActiveRow `json:"data"`
	TotalRecords *int                              `json:"total_records"`
}

type merchantDetailTrashedCacheResponseDB struct {
	Data         []*db.GetMerchantDetailsTrashedRow `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type merchantDetailQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantDetailQueryCache(store *cache.CacheStore) *merchantDetailQueryCache {
	return &merchantDetailQueryCache{
		store: store,
	}
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, *int, bool) {
	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDetailCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantDetailsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantDetailsRow{}
	}

	key := fmt.Sprintf(merchantDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, *int, bool) {
	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDetailActiveCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantDetailsActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantDetailsActiveRow{}
	}

	key := fmt.Sprintf(merchantDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, *int, bool) {
	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDetailTrashedCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantDetailsTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantDetailsTrashedRow{}
	}

	key := fmt.Sprintf(merchantDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDetailTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantDetailQueryCache) GetCachedMerchantDetail(ctx context.Context, id int) (*db.GetMerchantDetailRow, bool) {
	key := fmt.Sprintf(merchantDetailByIdCacheKey, id)

	result, found := cache.GetFromCache[db.GetMerchantDetailRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantDetailQueryCache) SetCachedMerchantDetail(ctx context.Context, data *db.GetMerchantDetailRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDetailByIdCacheKey, data.MerchantDetailID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
