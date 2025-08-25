package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Role RoleHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Role: NewRoleHandleGrpc(deps.Service, deps.Logger),
	}
}
