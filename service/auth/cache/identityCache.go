package cache

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

var (
	keyIdentityRefreshToken = "identity:refresh_token:%s"
	keyIdentityUserInfo     = "identity:user_info:%s"
)

type identityCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewidentityCache(store *sharedcachehelpers.CacheStore) *identityCache {
	return &identityCache{store: store}
}

func (c *identityCache) SetRefreshToken(ctx context.Context, token string, expiration time.Duration) {
	key := keyIdentityRefreshToken
	key = fmt.Sprintf(key, token)

	sharedcachehelpers.SetToCache(ctx, c.store, key, &token, expiration)
}

func (c *identityCache) GetRefreshToken(ctx context.Context, token string) (string, bool) {
	key := keyIdentityRefreshToken
	key = fmt.Sprintf(key, token)

	result, found := sharedcachehelpers.GetFromCache[string](ctx, c.store, key)
	if !found || result == nil {
		return "", false
	}
	return *result, true
}

func (c *identityCache) DeleteRefreshToken(ctx context.Context, token string) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

func (c *identityCache) SetCachedUserInfo(ctx context.Context, user *db.GetUserByIDRow, expiration time.Duration) {
	if user == nil {
		return
	}

	key := fmt.Sprintf(keyIdentityUserInfo, user.UserID)

	sharedcachehelpers.SetToCache(ctx, c.store, key, user, expiration)
}

func (c *identityCache) GetCachedUserInfo(ctx context.Context, userId string) (*db.GetUserByIDRow, bool) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	return sharedcachehelpers.GetFromCache[db.GetUserByIDRow](ctx, c.store, key)
}

func (c *identityCache) DeleteCachedUserInfo(ctx context.Context, userId string) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}
