package handler

import "github.com/MamangRust/monolith-ecommerce-auth/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	Auth AuthHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Auth: NewAuthHandleGrpc(
			deps.Service,
		),
	}
}
