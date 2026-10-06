package handler

import (
	"github.com/MamangRust/monolith-ecommerce-category/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	CategoryQuery           pbcategory.CategoryQueryServiceServer
	CategoryCommand         pbcategory.CategoryCommandServiceServer
	CategoryStats           CategoryStatsHandler
	CategoryStatsById       CategoryStatsByIdHandler
	CategoryStatsByMerchant CategoryStatsByMerchantHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		CategoryQuery:           NewCategoryQueryHandler(deps.Service.CategoryQuery, deps.Logger),
		CategoryCommand:         NewCategoryCommandHandler(deps.Service.CategoryCommand, deps.Logger),
		CategoryStats:           NewCategoryStatsHandler(deps.Service.CategoryStats, deps.Logger),
		CategoryStatsById:       NewCategoryStatsByIdHandler(deps.Service.CategoryStatsById, deps.Logger),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantHandler(deps.Service.CategoryStatsByMerchant, deps.Logger),
	}
}
