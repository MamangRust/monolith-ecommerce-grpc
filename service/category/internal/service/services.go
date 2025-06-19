package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-category/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	CategoryQuery           CategoryQueryService
	CategoryCommand         CategoryCommandService
	CategoryStats           CategoryStatsService
	CategoryStatsById       CategoryStatsByIdService
	CategoryStatsByMerchant CategoryStatsByMerchantService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	categoryMapper := response_service.NewCategoryResponseMapper()

	return &Service{
		CategoryQuery:           NewCategoryQueryService(deps.Ctx, deps.ErrorHandler.CategoryQueryError, deps.Mencache.CategoryQueryCache, deps.Repositories.CategoryQuery, deps.Logger, categoryMapper),
		CategoryCommand:         NewCategoryCommandService(deps.Ctx, deps.Mencache.CategoryCommandCache, deps.ErrorHandler.CategoryCommandError, deps.Repositories.CategoryCommand, deps.Repositories.CategoryQuery, deps.Logger, categoryMapper),
		CategoryStats:           NewCategoryStatsService(deps.Ctx, deps.Mencache.CategoryStatsCache, deps.ErrorHandler.CategoryStatsByIdError, deps.Repositories.CategoryStats, deps.Logger, categoryMapper),
		CategoryStatsById:       NewCategoryStatsByIdService(deps.Ctx, deps.Mencache.CategoryStatsByIdCache, deps.ErrorHandler.CategoryStatsByIdError, deps.Repositories.CategoryStatsById, deps.Logger, categoryMapper),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantService(deps.Ctx, deps.Mencache.CategoryStatsByMerchantCache, deps.ErrorHandler.CategoryStatsByMerchantError, deps.Repositories.CategoryStatsByMerchant, deps.Logger, categoryMapper),
	}
}
