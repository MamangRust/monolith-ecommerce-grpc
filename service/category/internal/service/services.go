package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-category/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	CategoryQuery           CategoryQueryService
	CategoryCommand         CategoryCommandService
	CategoryStats           CategoryStatsService
	CategoryStatsById       CategoryStatsByIdService
	CategoryStatsByMerchant CategoryStatsByMerchantService
}

type Deps struct {
	Cache         *mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		CategoryQuery: NewCategoryQueryService(&CategoryQueryServiceDeps{
			Observability:           deps.Observability,
			Cache:                   deps.Cache.CategoryQueryCache,
			CategoryQueryRepository: deps.Repositories.CategoryQuery,
			Logger:                  deps.Logger,
		}),
		CategoryCommand: NewCategoryCommandService(&CategoryCommandServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.CategoryCommandCache,
			CategoryCommandRepository: deps.Repositories.
				CategoryCommand,
			CategoryQueryRepository: deps.Repositories.CategoryQuery,
			Logger:                  deps.Logger,
		}),
		CategoryStats: NewCategoryStatsService(&CategoryStatsServiceDeps{
			Observability:           deps.Observability,
			Cache:                   deps.Cache.CategoryStatsCache,
			CategoryStatsRepository: deps.Repositories.CategoryStats,
			Logger:                  deps.Logger,
		}),
		CategoryStatsById: NewCategoryStatsByIdService(&CategoryStatsByIdServiceDeps{
			Observability:               deps.Observability,
			Cache:                       deps.Cache.CategoryStatsByIdCache,
			CategoryStatsByIdRepository: deps.Repositories.CategoryStatsById,
			Logger:                      deps.Logger,
		}),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantService(&CategoryStatsByMerchantServiceDeps{
			Observability:                     deps.Observability,
			Cache:                             deps.Cache.CategoryStatsByMerchantCache,
			CategoryStatsByMerchantRepository: deps.Repositories.CategoryStatsByMerchant,
			Logger:                            deps.Logger,
		}),
	}
}
