package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/service"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *service.Service
}

type Handler struct {
	OrderItem OrderItemHandlerGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderItem: NewOrderItemHandleGrpc(deps.Service.OrderItemQuery, protomapper.NewOrderItemProtoMapper()),
	}
}
