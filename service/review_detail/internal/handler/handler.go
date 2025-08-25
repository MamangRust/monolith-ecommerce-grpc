package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ReviewDetail ReviewDetailHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ReviewDetail: NewReviewDetailHandleGrpc(deps.Service, deps.Logger),
	}
}
