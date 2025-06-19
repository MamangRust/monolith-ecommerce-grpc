package handler

import "github.com/MamangRust/monolith-ecommerce-grpc-user/internal/service"

type Deps struct {
	Service *service.Service
}

type Handler struct {
	User UserHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		User: NewUserHandleGrpc(deps.Service),
	}
}
