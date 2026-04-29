package cache

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type roleCachedResponseAll struct {
	Data         []*db.GetRolesRow `json:"data"`
	TotalRecords *int              `json:"total_records"`
}

type roleCachedResponseActive struct {
	Data         []*db.GetActiveRolesRow `json:"data"`
	TotalRecords *int                    `json:"total_records"`
}

type roleCachedResponseTrashed struct {
	Data         []*db.GetTrashedRolesRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type roleQueryCache struct {
	store *cache.CacheStore
}

func NewRoleQueryCache(store *cache.CacheStore) RoleQueryCache {
	return &roleQueryCache{store: store}
}

func (r *roleQueryCache) SetCachedRoles(ctx context.Context, req *requests.FindAllRole, data []*db.GetRolesRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetRolesRow{}
	}

	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseAll{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoles(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, *int, bool) {
	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponseAll](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *roleQueryCache) SetCachedRoleById(ctx context.Context, id int, data *db.Role) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByIdCacheKey, id)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleById(ctx context.Context, id int) (*db.Role, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	result, found := cache.GetFromCache[*db.Role](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *roleQueryCache) SetCachedRoleByName(ctx context.Context, name string, data *db.Role) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByNameCacheKey, name)
	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleByName(ctx context.Context, name string) (*db.Role, bool) {
	key := fmt.Sprintf(roleByNameCacheKey, name)

	result, found := cache.GetFromCache[*db.Role](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *roleQueryCache) SetCachedRoleByUserId(ctx context.Context, userId int, data []*db.Role) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByUserIdCacheKey, userId)
	cache.SetToCache(ctx, r.store, key, &data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleByUserId(ctx context.Context, userId int) ([]*db.Role, bool) {
	key := fmt.Sprintf(roleByUserIdCacheKey, userId)

	result, found := cache.GetFromCache[[]*db.Role](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (r *roleQueryCache) SetCachedRoleActive(ctx context.Context, req *requests.FindAllRole, data []*db.GetActiveRolesRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetActiveRolesRow{}
	}

	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleActive(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, *int, bool) {
	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponseActive](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (r *roleQueryCache) SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole, data []*db.GetTrashedRolesRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*db.GetTrashedRolesRow{}
	}

	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, r.store, key, payload, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, *int, bool) {
	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[roleCachedResponseTrashed](ctx, r.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}
