package cache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
)

// registerCache is a struct that implements the RegisterCache interface
type registerCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewRegisterCache creates a new instance of registerCache.
// It sets up the CacheStore for the register caching operations.
// The function takes a CacheStore as a parameter and returns a pointer to the initialized registerCache.
func NewRegisterCache(store *sharedcachehelpers.CacheStore) *registerCache {
	return &registerCache{store: store}
}

// SetVerificationCodeCache stores a verification code in cache with expiration duration.
// Parameters:
//   - ctx: The context for the cache operation
//   - email: The user's email address as cache key
//   - code: The verification code string to store
//   - expiration: Duration until the code expires from cache
func (c *registerCache) SetVerificationCodeCache(ctx context.Context, email string, code string, expiration time.Duration) {
	if code == "" {
		return
	}

	key := fmt.Sprintf(keyVerifyCode, email)

	sharedcachehelpers.SetToCache(ctx, c.store, key, &code, expiration)
}
