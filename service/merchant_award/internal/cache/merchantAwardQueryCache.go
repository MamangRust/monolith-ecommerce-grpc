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
	merchantAwardAllCacheKey     = "merchant_award:all:page:%d:pageSize:%d:search:%s"
	merchantAwardByIdCacheKey    = "merchant_award:id:%d"
	merchantAwardActiveCacheKey  = "merchant_award:active:page:%d:pageSize:%d:search:%s"
	merchantAwardTrashedCacheKey = "merchant_award:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type merchantAwardCacheResponseDB struct {
	Data         []*db.GetMerchantCertificationsAndAwardsRow `json:"data"`
	TotalRecords *int                                        `json:"total_records"`
}

type merchantAwardActiveCacheResponseDB struct {
	Data         []*db.GetMerchantCertificationsAndAwardsActiveRow `json:"data"`
	TotalRecords *int                                              `json:"total_records"`
}

type merchantAwardTrashedCacheResponseDB struct {
	Data         []*db.GetMerchantCertificationsAndAwardsTrashedRow `json:"data"`
	TotalRecords *int                                               `json:"total_records"`
}

type merchantAwardQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantAwardQueryCache(store *cache.CacheStore) *merchantAwardQueryCache {
	return &merchantAwardQueryCache{
		store: store,
	}
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsRow, *int, bool) {
	key := fmt.Sprintf(merchantAwardAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantAwardCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantCertificationsAndAwardsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantCertificationsAndAwardsRow{}
	}

	key := fmt.Sprintf(merchantAwardAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantAwardCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, *int, bool) {
	key := fmt.Sprintf(merchantAwardActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantAwardActiveCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantCertificationsAndAwardsActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantCertificationsAndAwardsActiveRow{}
	}

	key := fmt.Sprintf(merchantAwardActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantAwardActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, *int, bool) {
	key := fmt.Sprintf(merchantAwardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantAwardTrashedCacheResponseDB](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantCertificationsAndAwardsTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantCertificationsAndAwardsTrashedRow{}
	}

	key := fmt.Sprintf(merchantAwardTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantAwardTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantAwardQueryCache) GetCachedMerchantAward(ctx context.Context, id int) (*db.GetMerchantCertificationOrAwardRow, bool) {
	key := fmt.Sprintf(merchantAwardByIdCacheKey, id)

	result, found := cache.GetFromCache[db.GetMerchantCertificationOrAwardRow](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (m *merchantAwardQueryCache) SetCachedMerchantAward(ctx context.Context, data *db.GetMerchantCertificationOrAwardRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantAwardByIdCacheKey, data.MerchantCertificationID)
	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}
