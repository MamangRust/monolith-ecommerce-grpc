package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type UserQueryCache interface {
	GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, bool)
	SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponse, total *int)

	GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool)
	SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int)

	GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool)
	SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int)

	GetCachedUserCache(ctx context.Context, id int) (*response.UserResponse, bool)
	SetCachedUserCache(ctx context.Context, data *response.UserResponse)
}

type UserCommandCache interface {
	DeleteUserCache(ctx context.Context, id int)
}
