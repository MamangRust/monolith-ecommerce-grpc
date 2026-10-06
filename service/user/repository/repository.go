package repository

import (
	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	pbuserrole "github.com/MamangRust/monolith-ecommerce-pb/user_role"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	roleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user_role"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type GuardOptions struct {
	UserRole []adapter.GuardOption
	Role     []adapter.GuardOption
}

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleRepository
	UserRole    UserRoleRepository
}

type Deps struct {
	Db              *db.Queries
	RoleQueryClient pbrole.RoleQueryServiceClient
	UserRoleClient  pbuserrole.UserRoleServiceClient
	Guard           GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(deps.Db),
		UserQuery:   NewUserQueryRepository(deps.Db),
		UserRole:    userroleadapter.New(deps.UserRoleClient, deps.Guard.UserRole...),
		Role:        roleadapter.New(deps.RoleQueryClient, deps.Guard.Role...),
	}
}
