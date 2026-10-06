package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-review/service"

	pbreview "github.com/MamangRust/monolith-ecommerce-pb/review"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ReviewQuery   pbreview.ReviewQueryServiceServer
	ReviewCommand pbreview.ReviewCommandServiceServer
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
