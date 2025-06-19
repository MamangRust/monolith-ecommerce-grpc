package mencache

import (
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

func (r *reviewDetailQueryCache) GetReviewDetailAllCache(req *requests.FindAllReview) ([]*response.ReviewDetailsResponse, *int, bool) {
	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewDetailCacheResponse](r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailAllCache(req *requests.FindAllReview, data []*response.ReviewDetailsResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.ReviewDetailsResponse{}
	}

	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponse{Data: data, Total: total}
	SetToCache(r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetRevieDetailActiveCache(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewDetailCacheResponseDeleteAt](r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailActiveCache(req *requests.FindAllReview, data []*response.ReviewDetailsResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.ReviewDetailsResponseDeleteAt{}
	}

	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponseDeleteAt{Data: data, Total: total}
	SetToCache(r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetReviewDetailTrashedCache(req *requests.FindAllReview) ([]*response.ReviewDetailsResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := GetFromCache[reviewDetailCacheResponseDeleteAt](r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.Total, true
}

func (r *reviewDetailQueryCache) SetReviewDetailTrashedCache(req *requests.FindAllReview, data []*response.ReviewDetailsResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.ReviewDetailsResponseDeleteAt{}
	}

	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponseDeleteAt{Data: data, Total: total}
	SetToCache(r.store, key, payload, ttlDefault)
}

func (r *reviewDetailQueryCache) GetCachedReviewDetailCache(review_id int) (*response.ReviewDetailsResponse, bool) {
	key := fmt.Sprintf(reviewByIdCacheKey, review_id)
	result, found := GetFromCache[*response.ReviewDetailsResponse](r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *reviewDetailQueryCache) SetCachedReviewDetailCache(data *response.ReviewDetailsResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(reviewByIdCacheKey, data.ID)
	SetToCache(r.store, key, data, ttlDefault)
}
