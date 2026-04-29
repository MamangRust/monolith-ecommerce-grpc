package cache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

var (
	keyPasswordResetToken = "password_reset:token:%s"

	keyVerifyCode = "register:verify_code:%s"
)

type passwordResetCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewPasswordResetCache initializes and returns a new instance of passwordResetCache.
// It sets up the CacheStore for the password reset caching operations.
// The function takes a CacheStore as a parameter and returns a pointer to the initialized passwordResetCache.
func NewPasswordResetCache(store *sharedcachehelpers.CacheStore) *passwordResetCache {
	return &passwordResetCache{store: store}
}

// SetResetTokenCache stores a password reset token associated with a user ID.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - token: the reset token string
//   - userID: the ID of the user who requested the reset
//   - expiration: the duration until the token expires
func (c *passwordResetCache) SetResetTokenCache(ctx context.Context, token string, userID int, expiration time.Duration) {
	key := fmt.Sprintf(keyPasswordResetToken, userID)

	sharedcachehelpers.SetToCache(ctx, c.store, key, &userID, expiration)
}

// GetResetTokenCache retrieves a user ID associated with a given reset token.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - token: the reset token to look up
//
// Returns:
//   - The associated user ID and a boolean indicating if the token was found.
func (c *passwordResetCache) GetResetTokenCache(ctx context.Context, token string) (int, bool) {
	key := fmt.Sprintf(keyPasswordResetToken, token)

	result, found := sharedcachehelpers.GetFromCache[int](ctx, c.store, key)

	if !found || result == nil {
		return 0, false
	}
	return *result, true
}

// DeleteResetTokenCache removes a reset token from the cache.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - token: the reset token to delete
func (c *passwordResetCache) DeleteResetTokenCache(ctx context.Context, token string) {
	key := fmt.Sprintf(keyPasswordResetToken, token)

	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

// DeleteVerificationCodeCache deletes the verification code associated with an email.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - email: the email address whose verification code should be deleted
func (c *passwordResetCache) DeleteVerificationCodeCache(ctx context.Context, email string) {
	key := fmt.Sprintf(keyVerifyCode, email)

	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}
