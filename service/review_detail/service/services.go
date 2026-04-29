package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	ReviewDetailQuery   ReviewDetailQueryService
	ReviewDetailCommand ReviewDetailCommandService
}

type Deps struct {
	Observability observability.TraceLoggerObservability
	Cache         *cache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	return &Service{
		ReviewDetailQuery: NewReviewDetailQueryService(&ReviewDetailQueryServiceDeps{
			Observability:          deps.Observability,
			Cache:                  deps.Cache.ReviewDetailQuery,
			ReviewDetailRepository: deps.Repositories.ReviewDetailQuery,
			Logger:                 deps.Logger,
		}),
		ReviewDetailCommand: NewReviewDetailCommandService(&ReviewDetailCommandServiceDeps{
			Observability:               deps.Observability,
			Cache:                       deps.Cache.ReviewDetailCommand,
			ReviewDetailRepository:      deps.Repositories.ReviewDetailCommand,
			ReviewDetailQueryRepository: deps.Repositories.ReviewDetailQuery,
			Logger:                      deps.Logger,
		}),
	}
}
