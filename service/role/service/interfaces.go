package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type RoleQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, *int, error)
	FindByID(ctx context.Context, role_id int) (*db.Role, error)
	FindByName(ctx context.Context, name string) (*db.Role, error)
	FindByUserId(ctx context.Context, id int) ([]*db.Role, error)
}

type RoleCommandService interface {
	Create(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error)
	Update(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error)
	Trash(ctx context.Context, role_id int) (*db.Role, error)
	Restore(ctx context.Context, role_id int) (*db.Role, error)
	DeletePermanent(ctx context.Context, role_id int) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)

	AssignRoleToUser(ctx context.Context, request *requests.CreateUserRoleRequest) (*db.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, request *requests.RemoveUserRoleRequest) error
}
