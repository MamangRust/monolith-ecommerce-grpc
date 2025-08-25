package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-product/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Product ProductHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Product: NewProductHandleGrpc(deps.Service, deps.Logger),
	}
}
