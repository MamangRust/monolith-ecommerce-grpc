package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type RoleQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, error)
	FindActive(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, error)
	FindByID(ctx context.Context, role_id int) (*db.Role, error)
	FindByName(ctx context.Context, name string) (*db.Role, error)
	FindByUserId(ctx context.Context, user_id int) ([]*db.Role, error)
}

type RoleCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error)
	Update(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error)
	Trash(ctx context.Context, role_id int) (*db.Role, error)
	Restore(ctx context.Context, role_id int) (*db.Role, error)
	DeletePermanent(ctx context.Context, role_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type UserRoleRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}
