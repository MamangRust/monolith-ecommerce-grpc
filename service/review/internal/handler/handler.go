package handler

import "github.com/MamangRust/monolith-ecommerce-grpc-review/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	Review ReviewHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Review: NewReviewHandleGrpc(deps.Service),
	}
}
