package handler

import "github.com/MamangRust/monolith-point-of-sale-grpc-order/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	Order OrderHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Order: NewOrderHandleGrpc(deps.Service),
	}
}
