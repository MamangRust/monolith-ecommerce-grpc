package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Review ReviewHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Review: NewReviewHandleGrpc(deps.Service, deps.Logger),
	}
}
