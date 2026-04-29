package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleRepository
}

func NewRepositories(DB *db.Queries, roleClient pb.RoleQueryServiceClient) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(DB),
		UserQuery:   NewUserQueryRepository(DB),
		Role:        NewRoleRepository(roleClient),
	}
}
