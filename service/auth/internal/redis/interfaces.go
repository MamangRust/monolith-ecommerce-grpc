package mencache

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type IdentityCache interface {
	SetRefreshToken(ctx context.Context, token string, expiration time.Duration)
	GetRefreshToken(ctx context.Context, token string) (string, bool)
	DeleteRefreshToken(ctx context.Context, token string)
	SetCachedUserInfo(ctx context.Context, user *response.UserResponse, expiration time.Duration)
	GetCachedUserInfo(ctx context.Context, userId string) (*response.UserResponse, bool)
	DeleteCachedUserInfo(ctx context.Context, userId string)
}

type LoginCache interface {
	SetCachedLogin(ctx context.Context, email string, data *response.TokenResponse, expiration time.Duration)
	GetCachedLogin(ctx context.Context, email string) (*response.TokenResponse, bool)
}

type PasswordResetCache interface {
	SetResetTokenCache(ctx context.Context, token string, userID int, expiration time.Duration)
	GetResetTokenCache(ctx context.Context, token string) (int, bool)
	DeleteResetTokenCache(ctx context.Context, token string)
	DeleteVerificationCodeCache(ctx context.Context, email string)
}

type RegisterCache interface {
	SetVerificationCodeCache(ctx context.Context, email string, code string, expiration time.Duration)
}
