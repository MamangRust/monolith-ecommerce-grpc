package handler

import "github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/service"

type Deps struct {
	Service *service.Service
}

type Handler struct {
	Banner BannerHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Banner: NewBannerHandleGrpc(deps.Service),
	}
}
