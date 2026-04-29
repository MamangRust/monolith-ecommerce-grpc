package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	CartQuery   CartQueryHandler
	CartCommand CartCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		CartQuery:   NewCartQueryHandler(deps.Service.CartQuery, deps.Logger),
		CartCommand: NewCartCommandHandler(deps.Service.CartCommand, deps.Logger),
	}
}
