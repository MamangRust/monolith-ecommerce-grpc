package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	merchantDocumentAllCacheKey     = "merchant_document:all:page:%d:pageSize:%d:search:%s"
	merchantDocumentByIdCacheKey    = "merchant_document:id:%d"
	merchantDocumentActiveCacheKey  = "merchant_document:active:page:%d:pageSize:%d:search:%s"
	merchantDocumentTrashedCacheKey = "merchant_document:trashed:page:%d:pageSize:%d:search:%s"
)

type merchantDocumentQueryCachedResponse struct {
	Data         []*response.MerchantDocumentResponse `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantDocumentQueryCachedResponseDeleteAt struct {
	Data         []*response.MerchantDocumentResponseDeleteAt `json:"data"`
	TotalRecords *int                                         `json:"total_records"`
}

type merchantDocumentQueryCache struct {
	store *CacheStore
}

func NewMerchantDocumentQueryCache(store *CacheStore) *merchantDocumentQueryCache {
	return &merchantDocumentQueryCache{store: store}
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, bool) {
	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDocumentQueryCachedResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponse{}
	}

	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponse{Data: data, TotalRecords: total}

	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDocumentQueryCachedResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseDeleteAt{Data: data, TotalRecords: total}

	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDocumentQueryCachedResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseDeleteAt{Data: data, TotalRecords: total}

	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocument(ctx context.Context, id int) (*response.MerchantDocumentResponse, bool) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)

	result, found := GetFromCache[*response.MerchantDocumentResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocument(ctx context.Context, data *response.MerchantDocumentResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDocumentByIdCacheKey, data.ID)

	SetToCache(ctx, s.store, key, data, ttlDefault)
}
