package mencache

import (
	"context"
	"fmt"
	"time"
)

type registerCache struct {
	store *CacheStore
}

func NewRegisterCache(store *CacheStore) *registerCache {
	return &registerCache{store: store}
}

func (c *registerCache) SetVerificationCodeCache(ctx context.Context, email string, code string, expiration time.Duration) {
	key := fmt.Sprintf(keyVerifyCode, email)

	SetToCache(ctx, c.store, key, &code, expiration)
}
