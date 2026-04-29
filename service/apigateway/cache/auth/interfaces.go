package auth_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type IdentityCache interface {
	SetRefreshToken(
		ctx context.Context,
		token string,
		response *response.ApiResponseRefreshToken,

	)
	GetRefreshToken(ctx context.Context, token string) (*response.ApiResponseRefreshToken, bool)
	DeleteRefreshToken(ctx context.Context, token string)

	SetCachedUserInfo(ctx context.Context, userId string, data *response.ApiResponseGetMe)
	GetCachedUserInfo(ctx context.Context, userId string) (*response.ApiResponseGetMe, bool)
	DeleteCachedUserInfo(ctx context.Context, userId string)
}

type LoginCache interface {
	GetCachedLogin(
		ctx context.Context,
		email string,
	) (*response.ApiResponseLogin, bool)
	SetCachedLogin(
		ctx context.Context,
		email string,
		data *response.ApiResponseLogin,
	)
}
