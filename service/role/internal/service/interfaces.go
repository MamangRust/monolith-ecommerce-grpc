package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type RoleQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponse, *int, *response.ErrorResponse)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(ctx context.Context, role_id int) (*response.RoleResponse, *response.ErrorResponse)
	FindByUserId(ctx context.Context, id int) ([]*response.RoleResponse, *response.ErrorResponse)
}

type RoleCommandService interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	TrashedRole(ctx context.Context, role_id int) (*response.RoleResponse, *response.ErrorResponse)
	RestoreRole(ctx context.Context, role_id int) (*response.RoleResponse, *response.ErrorResponse)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, *response.ErrorResponse)

	RestoreAllRole(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllRolePermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
