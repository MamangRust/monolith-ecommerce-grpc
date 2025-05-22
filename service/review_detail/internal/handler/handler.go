package handler

import "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	ReviewDetail ReviewDetailHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		ReviewDetail: NewReviewDetailHandleGrpc(deps.Service),
	}
}
