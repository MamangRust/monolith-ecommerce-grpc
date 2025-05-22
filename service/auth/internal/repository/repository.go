package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps Deps) *Repositories {
	return &Repositories{
		User:         NewUserRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRecordMapper),
		UserRole:     NewUserRoleRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRoleRecordMapper),
		RefreshToken: NewRefreshTokenRepository(deps.DB, deps.Ctx, deps.MapperRecord.RefreshTokenRecordMapper),
		Role:         NewRoleRepository(deps.DB, deps.Ctx, deps.MapperRecord.RoleRecordMapper),
		ResetToken:   NewResetTokenRepository(deps.DB, deps.Ctx, deps.MapperRecord.ResetTokenRecordMapper),
	}
}
