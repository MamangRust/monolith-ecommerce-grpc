package mencache

import (
	"context"
	"fmt"
	"time"
)

var (
	keyPasswordResetToken = "password_reset:token:%s"

	keyVerifyCode = "register:verify_code:%s"
)

type passwordResetCache struct {
	store *CacheStore
}

func NewPasswordResetCache(store *CacheStore) *passwordResetCache {
	return &passwordResetCache{store: store}
}

func (c *passwordResetCache) SetResetTokenCache(ctx context.Context, token string, userID int, expiration time.Duration) {
	key := fmt.Sprintf(keyPasswordResetToken, userID)

	SetToCache(ctx, c.store, key, &userID, expiration)
}

func (c *passwordResetCache) GetResetTokenCache(ctx context.Context, token string) (int, bool) {
	key := fmt.Sprintf(keyPasswordResetToken, token)

	result, found := GetFromCache[int](ctx, c.store, key)

	if !found || result == nil {
		return 0, false
	}

	return *result, true
}

func (c *passwordResetCache) DeleteResetTokenCache(ctx context.Context, token string) {
	key := fmt.Sprintf(keyPasswordResetToken, token)

	DeleteFromCache(ctx, c.store, key)
}

func (c *passwordResetCache) DeleteVerificationCodeCache(ctx context.Context, email string) {
	key := fmt.Sprintf(keyVerifyCode, email)

	DeleteFromCache(ctx, c.store, key)
}
