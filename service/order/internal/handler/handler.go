package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Order OrderHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Order: NewOrderHandleGrpc(deps.Service, deps.Logger),
	}
}
