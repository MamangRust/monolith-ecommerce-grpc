package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	reviewAllCacheKey     = "review:all:page:%d:pageSize:%d:search:%s"
	reviewByIdCacheKey    = "review:id:%d"
	reviewActiveCacheKey  = "review:active:page:%d:pageSize:%d:search:%s"
	reviewTrashedCacheKey = "review:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type reviewDetailCacheResponse struct {
	Data  []*response.ReviewDetailsResponse `json:"data"`
	Total *int                              `json:"total_records"`
}

type reviewDetailCacheResponseDeleteAt struct {
	Data  []*response.ReviewDetailsResponseDeleteAt `json:"data"`
	Total *int                                      `json:"total_records"`
}

type reviewDetailQueryCache struct {
	store *CacheStore
}

func NewReviewDetailQueryCache(store *CacheStore) *reviewDetailQueryCache {
	return &reviewDetailQueryCache{store: store}
}

func (r *reviewDetailQueryCache) GetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, bool) {
	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewDetailCacheResponse](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailAllCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewDetailsResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.ReviewDetailsResponse{}
	}

	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponse{Data: data, Total: total}
	SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewDetailCacheResponseDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailActiveCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewDetailsResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.ReviewDetailsResponseDeleteAt{}
	}

	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponseDeleteAt{Data: data, Total: total}
	SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[reviewDetailCacheResponseDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailTrashedCache(ctx context.Context, req *requests.FindAllReview, data []*response.ReviewDetailsResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.ReviewDetailsResponseDeleteAt{}
	}

	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponseDeleteAt{Data: data, Total: total}
	SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetCachedReviewDetailCache(ctx context.Context, review_id int) (*response.ReviewDetailsResponse, bool) {
	key := fmt.Sprintf(reviewByIdCacheKey, review_id)
	result, found := GetFromCache[*response.ReviewDetailsResponse](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *reviewDetailQueryCache) SetCachedReviewDetailCache(ctx context.Context, data *response.ReviewDetailsResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewByIdCacheKey, data.ID)
	SetToCache(ctx, r.store, key, data, ttlDefault)
}
