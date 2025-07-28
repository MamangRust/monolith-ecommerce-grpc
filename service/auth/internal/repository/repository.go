package repository

import (
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

func NewRepositories(DB *db.Queries) *Repositories {
	mapperUserRole := recordmapper.NewUserRoleRecordMapper()
	mapperUser := recordmapper.NewUserRecordMapper()
	mapperRefreshToken := recordmapper.NewRefreshTokenRecordMapper()
	mapperRole := recordmapper.NewRoleRecordMapper()
	mapperResetToken := recordmapper.NewResetTokenRecordMapper()

	return &Repositories{
		User:         NewUserRepository(DB, mapperUser),
		RefreshToken: NewRefreshTokenRepository(DB, mapperRefreshToken),
		UserRole:     NewUserRoleRepository(DB, mapperUserRole),
		Role:         NewRoleRepository(DB, mapperRole),
		ResetToken:   NewResetTokenRepository(DB, mapperResetToken),
	}
}
