package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	OrderItemQuery   OrderItemQueryHandler
	OrderItemCommand OrderItemCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderItemQuery:   NewOrderItemQueryHandler(deps.Service.OrderItemQuery, deps.Logger),
		OrderItemCommand: NewOrderItemCommandHandler(deps.Service.OrderItemCommand, deps.Logger),
	}
}
