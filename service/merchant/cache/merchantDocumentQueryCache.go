package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

const (
	merchantDocumentAllCacheKey     = "merchant_document:all:page:%d:pageSize:%d:search:%s"
	merchantDocumentByIdCacheKey    = "merchant_document:id:%d"
	merchantDocumentActiveCacheKey  = "merchant_document:active:page:%d:pageSize:%d:search:%s"
	merchantDocumentTrashedCacheKey = "merchant_document:trashed:page:%d:pageSize:%d:search:%s"
)

type merchantDocumentQueryCachedResponse struct {
	Data         []*db.GetMerchantDocumentsRow `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type merchantDocumentQueryCachedResponseDeleteAt struct {
	Data         []*db.GetActiveMerchantDocumentsRow `json:"data"`
	TotalRecords *int                                `json:"total_records"`
}

type merchantDocumentTrashedCachedResponse struct {
	Data         []*db.GetTrashedMerchantDocumentsRow `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantDocumentQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantDocumentQueryCache(store *cache.CacheStore) *merchantDocumentQueryCache {
	return &merchantDocumentQueryCache{store: store}
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDocumentQueryCachedResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetMerchantDocumentsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetMerchantDocumentsRow{}
	}

	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponse{Data: data, TotalRecords: total}

	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDocumentQueryCachedResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetActiveMerchantDocumentsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetActiveMerchantDocumentsRow{}
	}

	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseDeleteAt{Data: data, TotalRecords: total}

	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDocumentTrashedCachedResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetTrashedMerchantDocumentsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*db.GetTrashedMerchantDocumentsRow{}
	}

	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentTrashedCachedResponse{Data: data, TotalRecords: total}

	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocument(ctx context.Context, id int) (*db.GetMerchantDocumentRow, bool) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)

	result, found := cache.GetFromCache[db.GetMerchantDocumentRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocument(ctx context.Context, data *db.GetMerchantDocumentRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDocumentByIdCacheKey, data.DocumentID)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
