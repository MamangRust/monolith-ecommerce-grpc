package auth_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"fmt"
)

type identityCache struct {
	store *cache.CacheStore
}

func NewIdentityCache(store *cache.CacheStore) *identityCache {
	return &identityCache{store: store}
}

func (c *identityCache) SetRefreshToken(
	ctx context.Context,
	token string,
	response *response.ApiResponseRefreshToken,
) {
	if response == nil {
		return
	}

	key := fmt.Sprintf(keyIdentityRefreshToken, token)
	cache.SetToCache(ctx, c.store, key, response, ttlDefault)
}

func (c *identityCache) GetRefreshToken(ctx context.Context, token string) (*response.ApiResponseRefreshToken, bool) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)

	result, found := cache.GetFromCache[response.ApiResponseRefreshToken](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (c *identityCache) DeleteRefreshToken(
	ctx context.Context,
	token string,
) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)
	cache.DeleteFromCache(ctx, c.store, key)
}

func (c *identityCache) SetCachedUserInfo(ctx context.Context, userId string, data *response.ApiResponseGetMe) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	cache.SetToCache(ctx, c.store, key, data, ttlDefault)
}

func (c *identityCache) GetCachedUserInfo(ctx context.Context, userId string) (*response.ApiResponseGetMe, bool) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	result, found := cache.GetFromCache[response.ApiResponseGetMe](ctx, c.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (c *identityCache) DeleteCachedUserInfo(ctx context.Context, userId string) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)
	cache.DeleteFromCache(ctx, c.store, key)
}
