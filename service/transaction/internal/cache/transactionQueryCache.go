package cache

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

const (
	transactionAllCacheKey        = "transaction:all:page:%d:pageSize:%d:search:%s"
	transactionByIdCacheKey       = "transaction:id:%d"
	transactionByMerchantCacheKey = "transaction:merchant:%d:page:%d:pageSize:%d:search:%s"
	transactionActiveCacheKey     = "transaction:active:page:%d:pageSize:%d:search:%s"
	transactionTrashedCacheKey    = "transaction:trashed:page:%d:pageSize:%d:search:%s"
	transactionByOrderCacheKey    = "transaction:order:%d"
	ttlDefault                    = 5 * time.Minute
)

type transactionCacheResponseDB struct {
	Data         []*db.GetTransactionsRow `json:"data"`
	TotalRecords *int                     `json:"totalRecords"`
}

type transactionMechantCacheResponseDB struct {
	Data         []*db.GetTransactionByMerchantRow `json:"data"`
	TotalRecords *int                              `json:"totalRecords"`
}

type transactionActiveCacheResponseDB struct {
	Data         []*db.GetTransactionsActiveRow `json:"data"`
	TotalRecords *int                           `json:"totalRecords"`
}

type transactionTrashedCacheResponseDB struct {
	Data         []*db.GetTransactionsTrashedRow `json:"data"`
	TotalRecords *int                            `json:"totalRecords"`
}

type transactionQueryCache struct {
	store *cache.CacheStore
}

func NewTransactionQueryCache(store *cache.CacheStore) *transactionQueryCache {
	return &transactionQueryCache{store: store}
}

func (t *transactionQueryCache) GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, bool) {
	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[transactionCacheResponseDB](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTransactionsRow{}
	}

	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*db.GetTransactionByMerchantRow, *int, bool) {
	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[transactionMechantCacheResponseDB](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*db.GetTransactionByMerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTransactionByMerchantRow{}
	}

	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &transactionMechantCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, bool) {
	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[transactionActiveCacheResponseDB](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTransactionsActiveRow{}
	}

	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionActiveCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, bool) {
	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[transactionTrashedCacheResponseDB](ctx, t.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*db.GetTransactionsTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTransactionsTrashedRow{}
	}

	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionTrashedCacheResponseDB{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionCache(ctx context.Context, id int) (*db.GetTransactionByIDRow, bool) {
	key := fmt.Sprintf(transactionByIdCacheKey, id)
	result, found := cache.GetFromCache[db.GetTransactionByIDRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionCache(ctx context.Context, data *db.GetTransactionByIDRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionByIdCacheKey, data.TransactionID)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*db.GetTransactionByOrderIDRow, bool) {
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)
	result, found := cache.GetFromCache[db.GetTransactionByOrderIDRow](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *db.GetTransactionByOrderIDRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}
