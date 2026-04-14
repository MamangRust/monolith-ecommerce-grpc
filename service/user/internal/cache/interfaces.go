package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryCache interface {
	GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, *int, bool)
	SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*db.GetUsersRow, total *int)

	GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, *int, bool)
	SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*db.GetUsersActiveRow, total *int)

	GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, *int, bool)
	SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*db.GetUserTrashedRow, total *int)

	GetCachedUserCache(ctx context.Context, id int) (*db.GetUserByIDRow, bool)
	SetCachedUserCache(ctx context.Context, data *db.GetUserByIDRow)
}

type UserCommandCache interface {
	DeleteUserCache(ctx context.Context, id int)
}
