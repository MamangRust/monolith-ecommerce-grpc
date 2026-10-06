package repository

import (
	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	pbuserrole "github.com/MamangRust/monolith-ecommerce-pb/user_role"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"

	roleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/role"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user_role"
)

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

type GuardOptions struct {
	User     []adapter.GuardOption
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

var (
	_ UserRepository     = (*useradapter.Repository)(nil)
	_ RoleRepository     = (*roleadapter.Repository)(nil)
	_ UserRoleRepository = (*userroleadapter.Repository)(nil)
)

type Deps struct {
	Db                *db.Queries
	UserQueryClient   pbuser.UserQueryServiceClient
	UserCommandClient pbuser.UserCommandServiceClient
	RoleQueryClient   pbrole.RoleQueryServiceClient
	RoleCommandClient pbrole.RoleCommandServiceClient
	UserRoleClient    pbuserrole.UserRoleServiceClient
	Guard             GuardOptions
}

func NewRepositories(
	deps *Deps,
) *Repositories {
	return &Repositories{
		User:         useradapter.New(deps.UserQueryClient, deps.UserCommandClient, deps.Guard.User...),
		RefreshToken: NewRefreshTokenRepository(deps.Db),
		UserRole:     userroleadapter.New(deps.UserRoleClient, deps.Guard.UserRole...),
		Role:         roleadapter.New(deps.RoleQueryClient, deps.Guard.Role...),
		ResetToken:   NewResetTokenRepository(deps.Db),
	}
}
