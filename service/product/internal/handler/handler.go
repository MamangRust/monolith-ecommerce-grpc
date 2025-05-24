package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-product/internal/service"
)

type Deps struct {
	Service service.Service
}

type Handler struct {
	Product ProductHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Product: NewProductHandleGrpc(deps.Service),
	}
}
