package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	RoleCommand RoleCommandRepository
	RoleQuery   RoleQueryRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	roleMapper := recordmapper.NewRoleRecordMapper()

	return &Repositories{
		RoleCommand: NewRoleCommandRepository(deps.DB, deps.Ctx, roleMapper),
		RoleQuery:   NewRoleQueryRepository(deps.DB, deps.Ctx, roleMapper),
	}
}
