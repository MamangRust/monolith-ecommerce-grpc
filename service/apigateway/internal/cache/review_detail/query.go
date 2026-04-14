package reviewdetail_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

const (
	reviewDetailAllCacheKey         = "review_detail:all:page:%d:pageSize:%d:search:%s"
	reviewDetailByIdCacheKey        = "review_detail:id:%d"
	reviewDetailActiveCacheKey      = "review_detail:active:page:%d:pageSize:%d:search:%s"
	reviewDetailTrashedCacheKey     = "review_detail:trashed:page:%d:pageSize:%d:search:%s"
	reviewDetailByIdTrashedCacheKey = "review_detail:id_trashed:%d"

	ttlDefault = 5 * time.Minute
)

type reviewDetailQueryCache struct {
	store *cache.CacheStore
}

func NewReviewDetailQueryCache(store *cache.CacheStore) *reviewDetailQueryCache {
	return &reviewDetailQueryCache{store: store}
}

func (r *reviewDetailQueryCache) GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDetails, bool) {
	key := fmt.Sprintf(reviewDetailAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationReviewDetails](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDetails) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDetailsDeleteAt, bool) {
	key := fmt.Sprintf(reviewDetailActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationReviewDetailsDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDetailsDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDetailsDeleteAt, bool) {
	key := fmt.Sprintf(reviewDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationReviewDetailsDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDetailsDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewDetailQueryCache) GetCachedReviewDetailCache(ctx context.Context, reviewID int) (*response.ApiResponseReviewDetail, bool) {
	key := fmt.Sprintf(reviewDetailByIdCacheKey, reviewID)
	result, found := cache.GetFromCache[response.ApiResponseReviewDetail](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetCachedReviewDetailCache(ctx context.Context, data *response.ApiResponseReviewDetail) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailByIdCacheKey, data.Data.ID)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewDetailQueryCache) GetCachedReviewDetailTrashedCache(ctx context.Context, reviewID int) (*response.ApiResponseReviewDetailDeleteAt, bool) {
	key := fmt.Sprintf(reviewDetailByIdTrashedCacheKey, reviewID)
	result, found := cache.GetFromCache[response.ApiResponseReviewDetailDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewDetailQueryCache) SetCachedReviewDetailTrashedCache(ctx context.Context, data *response.ApiResponseReviewDetailDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewDetailByIdTrashedCacheKey, data.Data.ID)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}
