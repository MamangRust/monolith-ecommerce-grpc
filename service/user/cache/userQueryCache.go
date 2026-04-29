package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type userCacheResponseAll struct {
	Data         []*db.GetUsersRow `json:"data"`
	TotalRecords *int              `json:"total_records"`
}

type userCacheResponseActive struct {
	Data         []*db.GetUsersActiveRow `json:"data"`
	TotalRecords *int                    `json:"total_records"`
}

type userCacheResponseTrashed struct {
	Data         []*db.GetUserTrashedRow `json:"data"`
	TotalRecords *int                    `json:"total_records"`
}

type userQueryCache struct {
	store *cache.CacheStore
}

func NewUserQueryCache(store *cache.CacheStore) UserQueryCache {
	return &userQueryCache{store: store}
}

func (s *userQueryCache) GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, *int, bool) {
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponseAll](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*db.GetUsersRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetUsersRow{}
	}

	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseAll{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, *int, bool) {
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponseActive](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*db.GetUsersActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetUsersActiveRow{}
	}

	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, *int, bool) {
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponseTrashed](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*db.GetUserTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetUserTrashedRow{}
	}

	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserCache(ctx context.Context, id int) (*db.GetUserByIDRow, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	result, found := cache.GetFromCache[*db.GetUserByIDRow](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *userQueryCache) SetCachedUserCache(ctx context.Context, data *db.GetUserByIDRow) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userByIdCacheKey, data.UserID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
