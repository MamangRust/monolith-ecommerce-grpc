package cache

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/cache.go
type IdentityCache interface {
	SetRefreshToken(ctx context.Context, token string, expiration time.Duration)
	GetRefreshToken(ctx context.Context, token string) (string, bool)
	DeleteRefreshToken(ctx context.Context, token string)

	SetCachedUserInfo(ctx context.Context, user *db.GetUserByIDRow, expiration time.Duration)

	GetCachedUserInfo(ctx context.Context, userId string) (*db.GetUserByIDRow, bool)

	DeleteCachedUserInfo(ctx context.Context, userId string)
}

type LoginCache interface {
	SetCachedLogin(ctx context.Context, email string, data *response.TokenResponse, expiration time.Duration)
	GetCachedLogin(ctx context.Context, email string) (*response.TokenResponse, bool)
	IncrementFailedLogin(ctx context.Context, email string) (int, error)
	ResetFailedLogin(ctx context.Context, email string) error
	IsAccountLocked(ctx context.Context, email string) (bool, error)
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
