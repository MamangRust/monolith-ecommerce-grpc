package mencache

import (
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

func (s *merchantDocumentQueryCache) GetCachedMerchantDocuments(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, bool) {
	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDocumentQueryCachedResponse](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocuments(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponse{}
	}

	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponse{Data: data, TotalRecords: total}

	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsActive(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDocumentQueryCachedResponseDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsActive(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseDeleteAt{Data: data, TotalRecords: total}

	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsTrashed(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[merchantDocumentQueryCachedResponseDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsTrashed(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.MerchantDocumentResponseDeleteAt{}
	}

	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseDeleteAt{Data: data, TotalRecords: total}

	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocument(id int) (*response.MerchantDocumentResponse, bool) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)

	result, found := GetFromCache[*response.MerchantDocumentResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocument(data *response.MerchantDocumentResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDocumentByIdCacheKey, data.ID)

	SetToCache(s.store, key, data, ttlDefault)
}
