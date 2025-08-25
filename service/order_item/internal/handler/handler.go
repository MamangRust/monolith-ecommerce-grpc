package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	OrderItem OrderItemHandlerGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderItem: NewOrderItemHandleGrpc(deps.Service.OrderItemQuery, protomapper.NewOrderItemProtoMapper(), deps.Logger),
	}
}
