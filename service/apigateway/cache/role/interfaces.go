package role_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type RoleQueryCache interface {
	SetCachedRoles(ctx context.Context, req *requests.FindAllRole, data *response.ApiResponsePaginationRole)
	GetCachedRoles(ctx context.Context, req *requests.FindAllRole) (*response.ApiResponsePaginationRole, bool)

	SetCachedRoleById(ctx context.Context, id int, data *response.ApiResponseRole)
	GetCachedRoleById(ctx context.Context, id int) (*response.ApiResponseRole, bool)

	SetCachedRoleByUserId(ctx context.Context, userId int, data *response.ApiResponsesRole)
	GetCachedRoleByUserId(ctx context.Context, userId int) (*response.ApiResponsesRole, bool)

	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRole, data *response.ApiResponsePaginationRoleDeleteAt)
	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRole) (*response.ApiResponsePaginationRoleDeleteAt, bool)

	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole, data *response.ApiResponsePaginationRoleDeleteAt)
	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRole) (*response.ApiResponsePaginationRoleDeleteAt, bool)
}

type RoleCommandCache interface {
	DeleteCachedRole(ctx context.Context, id int)
}
