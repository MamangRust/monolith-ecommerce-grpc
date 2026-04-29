package role_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

type roleQueryCache struct {
	store *cache.CacheStore
}

func NewRoleQueryCache(store *cache.CacheStore) RoleQueryCache {
	return &roleQueryCache{store: store}
}

func (r *roleQueryCache) SetCachedRoles(ctx context.Context, req *requests.FindAllRole, data *response.ApiResponsePaginationRole) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoles(ctx context.Context, req *requests.FindAllRole) (*response.ApiResponsePaginationRole, bool) {
	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationRole](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *roleQueryCache) SetCachedRoleById(ctx context.Context, id int, data *response.ApiResponseRole) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByIdCacheKey, id)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleById(ctx context.Context, id int) (*response.ApiResponseRole, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	result, found := cache.GetFromCache[response.ApiResponseRole](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *roleQueryCache) SetCachedRoleByUserId(ctx context.Context, userId int, data *response.ApiResponsesRole) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleByUserIdCacheKey, userId)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleByUserId(ctx context.Context, userId int) (*response.ApiResponsesRole, bool) {
	key := fmt.Sprintf(roleByUserIdCacheKey, userId)

	result, found := cache.GetFromCache[response.ApiResponsesRole](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *roleQueryCache) SetCachedRoleActive(ctx context.Context, req *requests.FindAllRole, data *response.ApiResponsePaginationRoleDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleActive(ctx context.Context, req *requests.FindAllRole) (*response.ApiResponsePaginationRoleDeleteAt, bool) {
	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationRoleDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (r *roleQueryCache) SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole, data *response.ApiResponsePaginationRoleDeleteAt) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	cache.SetToCache(ctx, r.store, key, data, ttlDefault)
}

func (r *roleQueryCache) GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole) (*response.ApiResponsePaginationRoleDeleteAt, bool) {
	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[response.ApiResponsePaginationRoleDeleteAt](ctx, r.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}
