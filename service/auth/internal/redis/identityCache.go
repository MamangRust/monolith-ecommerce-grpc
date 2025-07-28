package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

var (
	keyIdentityRefreshToken = "identity:refresh_token:%s"
	keyIdentityUserInfo     = "identity:user_info:%s"
)

type identityCache struct {
	store *CacheStore
}

func NewidentityCache(store *CacheStore) *identityCache {
	return &identityCache{store: store}
}

func (c *identityCache) SetRefreshToken(ctx context.Context, token string, expiration time.Duration) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)

	SetToCache(ctx, c.store, key, &token, expiration)
}

func (c *identityCache) GetRefreshToken(ctx context.Context, token string) (string, bool) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)

	result, found := GetFromCache[string](ctx, c.store, key)

	if !found || result == nil {
		return "", false
	}

	return *result, true
}

func (c *identityCache) DeleteRefreshToken(ctx context.Context, token string) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)
	DeleteFromCache(ctx, c.store, key)
}

func (c *identityCache) SetCachedUserInfo(ctx context.Context, user *response.UserResponse, expiration time.Duration) {
	if user == nil {
		return
	}

	key := fmt.Sprintf(keyIdentityUserInfo, user.ID)

	SetToCache(ctx, c.store, key, user, expiration)
}

func (c *identityCache) GetCachedUserInfo(ctx context.Context, userId string) (*response.UserResponse, bool) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	return GetFromCache[response.UserResponse](ctx, c.store, key)
}

func (c *identityCache) DeleteCachedUserInfo(ctx context.Context, userId string) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	DeleteFromCache(ctx, c.store, key)
}
