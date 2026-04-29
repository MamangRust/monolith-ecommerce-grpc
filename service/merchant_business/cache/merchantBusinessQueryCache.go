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
	merchantBusinessAllCacheKey     = "merchant_business:all:page:%d:pageSize:%d:search:%s"
	merchantBusinessByIdCacheKey    = "merchant_business:id:%d"
	merchantBusinessActiveCacheKey  = "merchant_business:active:page:%d:pageSize:%d:search:%s"
	merchantBusinessTrashedCacheKey = "merchant_business:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantBusinessCacheResponseDB struct {
	Data         []*db.GetMerchantsBusinessInformationRow `json:"data"`
	TotalRecords *int                                     `json:"total_records"`
}

type merchantBusinessActiveCacheResponseDB struct {
	Data         []*db.GetMerchantsBusinessInformationActiveRow `json:"data"`
	TotalRecords *int                                           `json:"total_records"`
}

type merchantBusinessTrashedCacheResponseDB struct {
	Data         []*db.GetMerchantsBusinessInformationTrashedRow `json:"data"`
	TotalRecords *int                                            `json:"total_records"`
}

type merchantBusinessQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantBusinessQueryCache(store *cache.CacheStore) *merchantBusinessQueryCache {
	return &merchantBusinessQueryCache{
		store: store,
	}
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationRow, *int, bool) {
	key := fmt.Sprintf(merchantBusinessAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantBusinessCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsBusinessInformationRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantsBusinessInformationRow{}
	}

	key := fmt.Sprintf(merchantBusinessAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationActiveRow, *int, bool) {
	key := fmt.Sprintf(merchantBusinessActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantBusinessActiveCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsBusinessInformationActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantsBusinessInformationActiveRow{}
	}

	key := fmt.Sprintf(merchantBusinessActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationTrashedRow, *int, bool) {
	key := fmt.Sprintf(merchantBusinessTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantBusinessTrashedCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsBusinessInformationTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantsBusinessInformationTrashedRow{}
	}

	key := fmt.Sprintf(merchantBusinessTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantBusinessTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantBusinessQueryCache) GetCachedMerchantBusiness(ctx context.Context, id int) (*db.GetMerchantBusinessInformationRow, bool) {
	key := fmt.Sprintf(merchantBusinessByIdCacheKey, id)

	result, found := cache.GetFromCache[db.GetMerchantBusinessInformationRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantBusinessQueryCache) SetCachedMerchantBusiness(ctx context.Context, data *db.GetMerchantBusinessInformationRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantBusinessByIdCacheKey, data.MerchantBusinessInfoID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
