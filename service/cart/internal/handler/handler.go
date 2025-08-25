package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Cart CartHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Cart: NewCartHandleGrpc(deps.Service, deps.Logger),
	}
}
