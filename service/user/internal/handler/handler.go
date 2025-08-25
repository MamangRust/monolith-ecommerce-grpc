package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-user/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	User UserHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		User: NewUserHandleGrpc(deps.Service, deps.Logger),
	}
}
