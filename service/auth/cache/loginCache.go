package cache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

var (
	keylogin         = "auth:login:%s"
	keyFailedLogin   = "auth:login:failed:%s"
	keyAccountLocked = "auth:login:locked:%s"
)

const (
	maxFailedAttempts = 5
	lockoutDuration   = 15 * time.Minute
)

type loginCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewLoginCache creates and returns a new instance of loginCache.
// It initializes the cache with the provided CacheStore, which is used
// for storing and retrieving cached login-related data.
func NewLoginCache(store *sharedcachehelpers.CacheStore) *loginCache {
	return &loginCache{store: store}
}

// GetCachedLogin retrieves a cached token response by email.
func (s *loginCache) GetCachedLogin(ctx context.Context, email string) (*response.TokenResponse, bool) {
	key := fmt.Sprintf(keylogin, email)
	result, found := sharedcachehelpers.GetFromCache[*response.TokenResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetCachedLogin stores login token response in cache with expiration.
func (s *loginCache) SetCachedLogin(ctx context.Context, email string, data *response.TokenResponse, expiration time.Duration) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(keylogin, email)
	sharedcachehelpers.SetToCache(ctx, s.store, key, data, expiration)
}

func (s *loginCache) IncrementFailedLogin(ctx context.Context, email string) (int, error) {
	key := fmt.Sprintf(keyFailedLogin, email)

	// Increment the failed login counter in Redis
	val, err := s.store.Redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment failed login counter: %w", err)
	}

	// Set expiration for the counter if it's new (or refresh it)
	if val == 1 {
		s.store.Redis.Expire(ctx, key, 1*time.Hour)
	}

	// If threshold reached, lock the account
	if val >= maxFailedAttempts {
		lockKey := fmt.Sprintf(keyAccountLocked, email)
		s.store.Redis.Set(ctx, lockKey, true, lockoutDuration)
	}

	return int(val), nil
}

func (s *loginCache) ResetFailedLogin(ctx context.Context, email string) error {
	key := fmt.Sprintf(keyFailedLogin, email)
	lockKey := fmt.Sprintf(keyAccountLocked, email)

	err := s.store.Redis.Del(ctx, key, lockKey).Err()
	if err != nil {
		return fmt.Errorf("failed to reset failed login data: %w", err)
	}

	return nil
}

func (s *loginCache) IsAccountLocked(ctx context.Context, email string) (bool, error) {
	lockKey := fmt.Sprintf(keyAccountLocked, email)

	exists, err := s.store.Redis.Exists(ctx, lockKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check account lock status: %w", err)
	}

	return exists > 0, nil
}
