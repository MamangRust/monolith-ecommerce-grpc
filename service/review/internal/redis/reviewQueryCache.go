package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
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

type reviewCacheResponse struct {
	Data         []*response.ReviewResponse `json:"data"`
	TotalRecords *int                       `json:"total_records"`
}

type reviewCacheResponseDeleteAt struct {
	Data         []*response.ReviewResponseDeleteAt `json:"data"`
	TotalRecords *int                               `json:"total_records"`
}

type reviewDetailCacheResponse struct {
	Data         []*response.ReviewsDetailResponse `json:"data"`
	TotalRecords *int                              `json:"total_records"`
}

type reviewQueryCache struct {
	store *CacheStore
}

func NewReviewQueryCache(store *CacheStore) *reviewQueryCache {
	return &reviewQueryCache{store: store}
}

func (r *reviewQueryCache) GetReviewAllCache(req *requests.FindAllReview) ([]*response.ReviewResponse, *int, bool) {
	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewCacheResponse](r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewAllCache(req *requests.FindAllReview, data []*response.ReviewResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ReviewResponse{}
	}

	key := fmt.Sprintf(reviewAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewCacheResponse{Data: data, TotalRecords: total}
	SetToCache(r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewByProductCache(req *requests.FindAllReviewByProduct) ([]*response.ReviewsDetailResponse, *int, bool) {
	key := fmt.Sprintf(reviewProductCacheKey, req.ProductID, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewDetailCacheResponse](r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewByProductCache(req *requests.FindAllReviewByProduct, data []*response.ReviewsDetailResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ReviewsDetailResponse{}
	}

	key := fmt.Sprintf(reviewProductCacheKey, req.ProductID, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponse{Data: data, TotalRecords: total}
	SetToCache(r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewByMerchantCache(req *requests.FindAllReviewByMerchant) ([]*response.ReviewsDetailResponse, *int, bool) {
	key := fmt.Sprintf(reviewMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewDetailCacheResponse](r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewByMerchantCache(req *requests.FindAllReviewByMerchant, data []*response.ReviewsDetailResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ReviewsDetailResponse{}
	}

	key := fmt.Sprintf(reviewMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &reviewDetailCacheResponse{Data: data, TotalRecords: total}
	SetToCache(r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewActiveCache(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewCacheResponseDeleteAt](r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewActiveCache(req *requests.FindAllReview, data []*response.ReviewResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ReviewResponseDeleteAt{}
	}

	key := fmt.Sprintf(reviewActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(r.store, key, payload, ttlDefault)
}

func (r *reviewQueryCache) GetReviewTrashedCache(req *requests.FindAllReview) ([]*response.ReviewResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[reviewCacheResponseDeleteAt](r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *reviewQueryCache) SetReviewTrashedCache(req *requests.FindAllReview, data []*response.ReviewResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.ReviewResponseDeleteAt{}
	}

	key := fmt.Sprintf(reviewTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &reviewCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(r.store, key, payload, ttlDefault)
}
