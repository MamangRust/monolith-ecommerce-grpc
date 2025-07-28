package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	roleAllCacheKey     = "role:all:page:%d:pageSize:%d:search:%s"
	roleByIdCacheKey    = "role:id:%d"
	roleActiveCacheKey  = "role:active:page:%d:pageSize:%d:search:%s"
	roleTrashedCacheKey = "role:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type roleCachedResponse struct {
	Data         []*response.RoleResponse `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type roleCachedResponseDeleteAt struct {
	Data         []*response.RoleResponseDeleteAt `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type roleQueryCache struct {
	store *CacheStore
}

func NewRoleQueryCache(store *CacheStore) *roleQueryCache {
	return &roleQueryCache{store: store}
}

func (m *roleQueryCache) SetCachedRoles(ctx context.Context, req *requests.FindAllRole, data []*response.RoleResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.RoleResponse{}
	}

	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &roleCachedResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *roleQueryCache) SetCachedRoleById(ctx context.Context, data *response.RoleResponse) {
	if data == nil {
		data = &response.RoleResponse{}
	}

	key := fmt.Sprintf(roleByIdCacheKey, data.ID)
	SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *roleQueryCache) SetCachedRoleByUserId(ctx context.Context, userId int, data []*response.RoleResponse) {
	if data == nil {
		data = []*response.RoleResponse{}
	}

	key := fmt.Sprintf(roleByIdCacheKey, userId)

	SetToCache(ctx, m.store, key, &data, ttlDefault)
}

func (m *roleQueryCache) SetCachedRoleActive(ctx context.Context, req *requests.FindAllRole, data []*response.RoleResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.RoleResponseDeleteAt{}
	}

	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &roleCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *roleQueryCache) SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole, data []*response.RoleResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.RoleResponseDeleteAt{}
	}

	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *roleQueryCache) GetCachedRoles(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponse, *int, bool) {
	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[roleCachedResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *roleQueryCache) GetCachedRoleById(ctx context.Context, id int) (*response.RoleResponse, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	result, found := GetFromCache[*response.RoleResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *roleQueryCache) GetCachedRoleByUserId(ctx context.Context, userId int) ([]*response.RoleResponse, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, userId)

	result, found := GetFromCache[[]*response.RoleResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *roleQueryCache) GetCachedRoleActive(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[roleCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *roleQueryCache) GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[roleCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}
