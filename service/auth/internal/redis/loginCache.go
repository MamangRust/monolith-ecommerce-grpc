package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

var keylogin = "auth:login:%s"

type loginCache struct {
	store *CacheStore
}

func NewLoginCache(store *CacheStore) *loginCache {
	return &loginCache{store: store}
}

func (s *loginCache) GetCachedLogin(email string) (*response.TokenResponse, bool) {
	key := fmt.Sprintf(keylogin, email)

	result, found := GetFromCache[*response.TokenResponse](s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *loginCache) SetCachedLogin(email string, data *response.TokenResponse, expiration time.Duration) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(keylogin, email)

	SetToCache(s.store, key, data, expiration)
}
