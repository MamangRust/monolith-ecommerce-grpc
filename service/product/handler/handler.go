package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-product/service"
)

type Handler struct {
	ProductQuery   ProductQueryHandler
	ProductCommand ProductCommandHandler
}

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ProductQuery:   NewProductQueryHandler(deps.Service.ProductQuery, deps.Logger),
		ProductCommand: NewProductCommandHandler(deps.Service.ProductCommand, deps.Logger),
	}
}
