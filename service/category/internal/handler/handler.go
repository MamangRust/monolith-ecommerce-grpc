package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Category CategoryHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Category: NewCategoryHandleGrpc(deps.Service, deps.Logger),
	}
}
