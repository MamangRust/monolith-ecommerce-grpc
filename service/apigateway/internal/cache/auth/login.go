package auth_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
	"time"
)

var ttlDefault = 5 * time.Minute

type loginCache struct {
	store *cache.CacheStore
}

func NewLoginCache(store *cache.CacheStore) *loginCache {
	return &loginCache{store: store}
}

func (s *loginCache) GetCachedLogin(ctx context.Context, email string) (*response.ApiResponseLogin, bool) {
	key := fmt.Sprintf(keylogin, email)

	result, found := cache.GetFromCache[response.ApiResponseLogin](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *loginCache) SetCachedLogin(ctx context.Context, email string, data *response.ApiResponseLogin) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(keylogin, email)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
