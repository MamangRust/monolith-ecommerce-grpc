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
	reviewAllCacheKey      = "review:all:page:%d:pageSize:%d:search:%s"
	reviewProductCacheKey  = "review:product:%d:page:%d:pageSize:%d:search:%s"
	reviewMerchantCacheKey = "review:merchant:%d:page:%d:pageSize:%d:search:%s"

	reviewByIdCacheKey    = "review:id:%d"
	reviewActiveCacheKey  = "review:active:page:%d:pageSize:%d:search:%s"
	reviewTrashedCacheKey = "review:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type reviewCacheResponseDB struct {
	Data         []*db.GetReviewsRow `json:"data"`
	TotalRecords *int                `json:"total_records"`
}

type reviewByProductCacheResponseDB struct {
	Data         []*db.GetReviewByProductIdRow `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type reviewByMerchantCacheResponseDB struct {
	Data         []*db.GetReviewByMerchantIdRow `json:"data"`
	TotalRecords *int                           `json:"total_records"`
}

type reviewActiveCacheResponseDB struct {
	Data         []*db.GetReviewsActiveRow `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type reviewTrashedCacheResponseDB struct {
	Data         []*db.GetReviewsTrashedRow `json:"data"`
	TotalRecords *int                       `json:"total_records"`
}

type reviewQueryCache struct {
	store *cache.CacheStore
}

func NewReviewQueryCache(store *cache.CacheStore) *reviewQueryCache {
	return &reviewQueryCache{store: store}
}

func (r *reviewQueryCache) GetReviewAllCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, *int, bool) {
	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewAllCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetReviewsRow{}
	}

	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, *int, bool) {
	key := fmt.Sprintf(reviewProductCacheKey, req.ProductID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewByProductCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct, data []*db.GetReviewByProductIdRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetReviewByProductIdRow{}
	}

	key := fmt.Sprintf(reviewProductCacheKey, req.ProductID, req.Page, req.PageSize, req.Search)
	payload := &reviewByProductCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, *int, bool) {
	key := fmt.Sprintf(reviewMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewByMerchantCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant, data []*db.GetReviewByMerchantIdRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetReviewByMerchantIdRow{}
	}

	key := fmt.Sprintf(reviewMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &reviewByMerchantCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, *int, bool) {
	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewActiveCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewActiveCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewsActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetReviewsActiveRow{}
	}

	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, *int, bool) {
	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewTrashedCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewsTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetReviewsTrashedRow{}
	}

	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewByIdCache(ctx context.Context, id int) (*db.GetReviewByIDRow, bool) {
	key := fmt.Sprintf(reviewByIdCacheKey, id)
	result, found := cache.GetFromCache[db.GetReviewByIDRow](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewQueryCache) SetReviewByIdCache(ctx context.Context, data *db.GetReviewByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewByIdCacheKey, data.ReviewID)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}
