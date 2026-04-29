package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type RoleQueryCache interface {
	SetCachedRoles(ctx context.Context, req *requests.FindAllRole, data []*db.GetRolesRow, total *int)
	GetCachedRoles(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, *int, bool)

	GetCachedRoleById(ctx context.Context, id int) (*db.Role, bool)
	SetCachedRoleById(ctx context.Context, id int, data *db.Role)

	GetCachedRoleByName(ctx context.Context, name string) (*db.Role, bool)
	SetCachedRoleByName(ctx context.Context, name string, data *db.Role)

	GetCachedRoleByUserId(ctx context.Context, userId int) ([]*db.Role, bool)
	SetCachedRoleByUserId(ctx context.Context, userId int, data []*db.Role)

	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, *int, bool)
	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRole, data []*db.GetActiveRolesRow, total *int)

	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, *int, bool)
	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole, data []*db.GetTrashedRolesRow, total *int)
}

type RoleCommandCache interface {
	DeleteCachedRole(ctx context.Context, id int)
}
