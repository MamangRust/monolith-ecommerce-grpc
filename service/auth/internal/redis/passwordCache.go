package mencache

import (
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

func (c *passwordResetCache) SetResetTokenCache(token string, userID int, expiration time.Duration) {
	key := fmt.Sprintf(keyPasswordResetToken, userID)

	SetToCache(c.store, key, &userID, expiration)
}

func (c *passwordResetCache) GetResetTokenCache(token string) (int, bool) {
	key := fmt.Sprintf(keyPasswordResetToken, token)

	result, found := GetFromCache[int](c.store, key)

	if !found || result == nil {
		return 0, false
	}

	return *result, true
}

func (c *passwordResetCache) DeleteResetTokenCache(token string) {
	key := fmt.Sprintf(keyPasswordResetToken, token)

	DeleteFromCache(c.store, key)
}

func (c *passwordResetCache) DeleteVerificationCodeCache(email string) {
	key := fmt.Sprintf(keyVerifyCode, email)

	DeleteFromCache(c.store, key)
}
