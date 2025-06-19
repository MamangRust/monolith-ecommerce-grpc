package mencache

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type IdentityCache interface {
	SetRefreshToken(token string, expiration time.Duration)
	GetRefreshToken(token string) (string, bool)
	DeleteRefreshToken(token string)
	SetCachedUserInfo(user *response.UserResponse, expiration time.Duration)
	GetCachedUserInfo(userId string) (*response.UserResponse, bool)
	DeleteCachedUserInfo(userId string)
}

type LoginCache interface {
	SetCachedLogin(email string, data *response.TokenResponse, expiration time.Duration)
	GetCachedLogin(email string) (*response.TokenResponse, bool)
}

type PasswordResetCache interface {
	SetResetTokenCache(token string, userID int, expiration time.Duration)
	GetResetTokenCache(token string) (int, bool)
	DeleteResetTokenCache(token string)
	DeleteVerificationCodeCache(email string)
}

type RegisterCache interface {
	SetVerificationCodeCache(email string, code string, expiration time.Duration)
}
