package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	RoleQuery   RoleQueryService
	RoleCommand RoleCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	roleMapper := response_service.NewRoleResponseMapper()

	return &Service{
		RoleQuery:   NewRoleQueryService(deps.Ctx, deps.Repositories.RoleQuery, deps.Logger, roleMapper),
		RoleCommand: NewRoleCommandService(deps.Ctx, deps.Repositories.RoleCommand, deps.Logger, roleMapper),
	}
}
