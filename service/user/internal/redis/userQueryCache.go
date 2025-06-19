package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

const (
	userAllCacheKey     = "user:all:page:%d:pageSize:%d:search:%s"
	userByIdCacheKey    = "user:id:%d"
	userActiveCacheKey  = "user:active:page:%d:pageSize:%d:search:%s"
	userTrashedCacheKey = "user:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type userCacheResponse struct {
	Data         []*response.UserResponse `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type userCacheResponseDeleteAt struct {
	Data         []*response.UserResponseDeleteAt `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

type userQueryCache struct {
	store *CacheStore
}

func NewUserQueryCache(store *CacheStore) *userQueryCache {
	return &userQueryCache{store: store}
}

func (s *userQueryCache) GetCachedUsersCache(req *requests.FindAllUsers) ([]*response.UserResponse, *int, bool) {
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[userCacheResponse](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUsersCache(req *requests.FindAllUsers, data []*response.UserResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.UserResponse{}
	}

	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponse{Data: data, TotalRecords: total}
	SetToCache(s.store, key, payload, ttlDefault)

}

func (s *userQueryCache) GetCachedUserActiveCache(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[userCacheResponseDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserActiveCache(req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.UserResponseDeleteAt{}
	}

	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserTrashedCache(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[userCacheResponseDeleteAt](s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserTrashedCache(req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.UserResponseDeleteAt{}
	}

	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserCache(id int) (*response.UserResponse, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	result, found := GetFromCache[*response.UserResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true

}
func (s *userQueryCache) SetCachedUserCache(data *response.UserResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userByIdCacheKey, data.ID)
	SetToCache(s.store, key, data, ttlDefault)
}
