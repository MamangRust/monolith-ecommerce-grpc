package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	CategoryQuery           pb.CategoryQueryServiceServer
	CategoryCommand         pb.CategoryCommandServiceServer
	CategoryStats           CategoryStatsHandler
	CategoryStatsById       CategoryStatsByIdHandler
	CategoryStatsByMerchant CategoryStatsByMerchantHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		CategoryQuery:           NewCategoryQueryHandler(deps.Service.CategoryQuery, deps.Logger),
		CategoryCommand: NewCategoryCommandHandler(deps.Service.CategoryCommand, deps.Logger),
		CategoryStats:           NewCategoryStatsHandler(deps.Service, deps.Logger),
		CategoryStatsById:       NewCategoryStatsByIdHandler(deps.Service, deps.Logger),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantHandler(deps.Service, deps.Logger),
	}
}
