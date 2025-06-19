package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type RoleCommandCache interface {
	DeleteCachedRole(id int)
}

type RoleQueryCache interface {
	SetCachedRoles(req *requests.FindAllRole, data []*response.RoleResponse, total *int)
	SetCachedRoleById(data *response.RoleResponse)
	SetCachedRoleByUserId(userId int, data []*response.RoleResponse)
	SetCachedRoleActive(req *requests.FindAllRole, data []*response.RoleResponseDeleteAt, total *int)
	SetCachedRoleTrashed(req *requests.FindAllRole, data []*response.RoleResponseDeleteAt, total *int)

	GetCachedRoles(req *requests.FindAllRole) ([]*response.RoleResponse, *int, bool)
	GetCachedRoleByUserId(userId int) ([]*response.RoleResponse, bool)
	GetCachedRoleById(id int) (*response.RoleResponse, bool)
	GetCachedRoleActive(req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, bool)
	GetCachedRoleTrashed(req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, bool)
}
