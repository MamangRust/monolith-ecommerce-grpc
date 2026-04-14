package review_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
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

// Struktur pembungkus (reviewCacheResponseDB, dll.) sudah tidak diperlukan lagi
// karena tipe ApiResponse... sudah mencakup data dan paginasi.

type reviewQueryCache struct {
	store *cache.CacheStore
}

func NewReviewQueryCache(store *cache.CacheStore) *reviewQueryCache {
	return &reviewQueryCache{store: store}
}

func (r *reviewQueryCache) GetReviewAllCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReview, bool) {
	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationReview](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewQueryCache) SetReviewAllCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReview) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewQueryCache) GetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct) (*response.ApiResponsePaginationReview, bool) {
	key := fmt.Sprintf(reviewProductCacheKey, req.ProductID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationReview](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewQueryCache) SetReviewByProductCache(ctx context.Context, req *requests.FindAllReviewByProduct, data *response.ApiResponsePaginationReview) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewProductCacheKey, req.ProductID, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewQueryCache) GetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant) (*response.ApiResponsePaginationReviewsDetail, bool) {
	key := fmt.Sprintf(reviewMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationReviewsDetail](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewQueryCache) SetReviewByMerchantCache(ctx context.Context, req *requests.FindAllReviewByMerchant, data *response.ApiResponsePaginationReviewsDetail) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewQueryCache) GetReviewActiveCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDeleteAt, bool) {
	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationReviewDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewQueryCache) SetReviewActiveCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewQueryCache) GetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview) (*response.ApiResponsePaginationReviewDeleteAt, bool) {
	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[response.ApiResponsePaginationReviewDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewQueryCache) SetReviewTrashedCache(ctx context.Context, req *requests.FindAllReview, data *response.ApiResponsePaginationReviewDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *reviewQueryCache) GetReviewByIdCache(ctx context.Context, id int) (*response.ApiResponseReview, bool) {
	key := fmt.Sprintf(reviewByIdCacheKey, id)
	result, found := cache.GetFromCache[response.ApiResponseReview](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *reviewQueryCache) SetReviewByIdCache(ctx context.Context, data *response.ApiResponseReview) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewByIdCacheKey, data.Data.ID)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}
