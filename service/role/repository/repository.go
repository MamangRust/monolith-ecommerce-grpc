package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
)

type Repositories struct {
	RoleCommand RoleCommandRepository
	RoleQuery   RoleQueryRepository
	UserRole    UserRoleRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		RoleCommand: NewRoleCommandRepository(DB),
		RoleQuery:   NewRoleQueryRepository(DB),
		UserRole:    NewUserRoleRepository(DB),
	}
}
