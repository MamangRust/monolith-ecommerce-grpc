package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Banner BannerHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Banner: NewBannerHandleGrpc(deps.Service, deps.Logger),
	}
}
