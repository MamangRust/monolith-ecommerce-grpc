package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type RoleQueryRepository interface {
	FindAllRoles(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error)
	FindByActiveRole(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error)
	FindByTrashedRole(req *requests.FindAllRole) ([]*record.RoleRecord, *int, error)
	FindById(role_id int) (*record.RoleRecord, error)
	FindByName(name string) (*record.RoleRecord, error)
	FindByUserId(user_id int) ([]*record.RoleRecord, error)
}

type RoleCommandRepository interface {
	CreateRole(request *requests.CreateRoleRequest) (*record.RoleRecord, error)
	UpdateRole(request *requests.UpdateRoleRequest) (*record.RoleRecord, error)
	TrashedRole(role_id int) (*record.RoleRecord, error)
	RestoreRole(role_id int) (*record.RoleRecord, error)
	DeleteRolePermanent(role_id int) (bool, error)
	RestoreAllRole() (bool, error)
	DeleteAllRolePermanent() (bool, error)
}
