package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/service"
)

type Deps struct {
	Service service.Service
}

type Handler struct {
	Cart CartHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Cart: NewCartHandleGrpc(deps.Service),
	}
}
