package mencache

import (
	"fmt"
	"time"
)

type registerCache struct {
	store *CacheStore
}

func NewRegisterCache(store *CacheStore) *registerCache {
	return &registerCache{store: store}
}

func (c *registerCache) SetVerificationCodeCache(email string, code string, expiration time.Duration) {
	key := fmt.Sprintf(keyVerifyCode, email)

	SetToCache(c.store, key, &code, expiration)
}
