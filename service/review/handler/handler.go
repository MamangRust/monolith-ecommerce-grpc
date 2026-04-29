package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-review/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ReviewQuery   pb.ReviewQueryServiceServer
	ReviewCommand pb.ReviewCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ReviewQuery:   NewReviewQueryHandler(deps.Service.ReviewQuery, deps.Logger),
		ReviewCommand: NewReviewCommandHandler(deps.Service.ReviewCommand, deps.Logger),
	}
}

type reviewHandleGrpc struct {
	// Dummy struct for mapping receiver
}
