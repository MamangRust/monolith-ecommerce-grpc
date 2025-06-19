package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	transactionAllCacheKey  = "transaction:all:page:%d:pageSize:%d:search:%s"
	transactionByIdCacheKey = "transaction:id:%d"

	transactionByMerchantCacheKey = "transaction:merchant:%d:page:%d:pageSize:%d:search:%s"

	transactionActiveCacheKey  = "transaction:active:page:%d:pageSize:%d:search:%s"
	transactionTrashedCacheKey = "transaction:trashed:page:%d:pageSize:%d:search:%s"

	transactionByOrderCacheKey = "transaction:order:%d"

	ttlDefault = 5 * time.Minute
)

type transactionCacheResponse struct {
	Data         []*response.TransactionResponse `json:"data"`
	TotalRecords *int                            `json:"totalRecords"`
}

type transactionCacheResponseDeleteAt struct {
	Data         []*response.TransactionResponseDeleteAt `json:"data"`
	TotalRecords *int                                    `json:"totalRecords"`
}

type transactionQueryCache struct {
	store *CacheStore
}

func NewTransactionQueryCache(store *CacheStore) *transactionQueryCache {
	return &transactionQueryCache{store: store}
}

func (t *transactionQueryCache) GetCachedTransactionsCache(req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, bool) {
	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transactionCacheResponse](t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionsCache(req *requests.FindAllTransaction, data []*response.TransactionResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TransactionResponse{}
	}

	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCacheResponse{Data: data, TotalRecords: total}
	SetToCache(t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByMerchant(req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, bool) {
	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transactionCacheResponse](t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionByMerchant(req *requests.FindAllTransactionByMerchant, data []*response.TransactionResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TransactionResponse{}
	}

	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	payload := &transactionCacheResponse{Data: data, TotalRecords: total}
	SetToCache(t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionActiveCache(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transactionCacheResponseDeleteAt](t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionActiveCache(req *requests.FindAllTransaction, data []*response.TransactionResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TransactionResponseDeleteAt{}
	}

	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionTrashedCache(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[transactionCacheResponseDeleteAt](t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionTrashedCache(req *requests.FindAllTransaction, data []*response.TransactionResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.TransactionResponseDeleteAt{}
	}

	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &transactionCacheResponseDeleteAt{Data: data, TotalRecords: total}

	SetToCache(t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionCache(id int) (*response.TransactionResponse, bool) {
	key := fmt.Sprintf(transactionByIdCacheKey, id)

	result, found := GetFromCache[*response.TransactionResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionQueryCache) SetCachedTransactionCache(data *response.TransactionResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transactionByIdCacheKey, data.ID)
	SetToCache(t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByOrderId(orderID int) (*response.TransactionResponse, bool) {
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)

	result, found := GetFromCache[*response.TransactionResponse](t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionQueryCache) SetCachedTransactionByOrderId(orderID int, data *response.TransactionResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)
	SetToCache(t.store, key, data, ttlDefault)
}
