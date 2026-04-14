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
	reviewDetailAllCacheKey         = "review_detail:all:page:%d:pageSize:%d:search:%s"
	reviewDetailByIdCacheKey        = "review_detail:id:%d"
	reviewDetailActiveCacheKey      = "review_detail:active:page:%d:pageSize:%d:search:%s"
	reviewDetailTrashedCacheKey     = "review_detail:trashed:page:%d:pageSize:%d:search:%s"
	reviewDetailByIdTrashedCacheKey = "review_detail:id_trashed:%d"

	ttlDefault = 5 * time.Minute
)

type reviewDetailCacheResponseDB struct {
	Data  []*db.GetReviewDetailsRow `json:"data"`
	Total *int                      `json:"total_records"`
}

type reviewDetailActiveCacheResponseDB struct {
	Data  []*db.GetReviewDetailsActiveRow `json:"data"`
	Total *int                            `json:"total_records"`
}

type reviewDetailTrashedCacheResponseDB struct {
	Data  []*db.GetReviewDetailsTrashedRow `json:"data"`
	Total *int                             `json:"total_records"`
}

type reviewDetailQueryCache struct {
	store *cache.CacheStore
}

func NewReviewDetailQueryCache(store *cache.CacheStore) *reviewDetailQueryCache {
	return &reviewDetailQueryCache{store: store}
}

func (r *reviewDetailQueryCache) GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, *int, bool) {
	key := fmt.Sprintf(reviewDetailAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewDetailCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewDetailsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetReviewDetailsRow{}
	}

	key := fmt.Sprintf(reviewDetailAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, *int, bool) {
	key := fmt.Sprintf(reviewDetailActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewDetailActiveCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewDetailsActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetReviewDetailsActiveRow{}
	}

	key := fmt.Sprintf(reviewDetailActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailActiveCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, *int, bool) {
	key := fmt.Sprintf(reviewDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[reviewDetailTrashedCacheResponseDB](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*db.GetReviewDetailsTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetReviewDetailsTrashedRow{}
	}

	key := fmt.Sprintf(reviewDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailTrashedCacheResponseDB{Data: data, Total: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetCachedReviewDetailCache(ctx context.Context, reviewID int) (*db.GetReviewDetailRow, bool) {
	key := fmt.Sprintf(reviewDetailByIdCacheKey, reviewID)
	result, found := cache.GetFromCache[db.GetReviewDetailRow](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetCachedReviewDetailCache(ctx context.Context, data *db.GetReviewDetailRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailByIdCacheKey, data.ReviewDetailID)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewDetailQueryCache) GetCachedReviewDetailTrashedCache(ctx context.Context, reviewID int) (*db.ReviewDetail, bool) {
	key := fmt.Sprintf(reviewDetailByIdTrashedCacheKey, reviewID)
	result, found := cache.GetFromCache[db.ReviewDetail](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetCachedReviewDetailTrashedCache(ctx context.Context, data *db.ReviewDetail) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailByIdTrashedCacheKey, data.ReviewDetailID)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}
