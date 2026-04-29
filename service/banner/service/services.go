package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-banner/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	BannerQuery   BannerQueryService
	BannerCommand BannerCommandService
}

type Deps struct {
	Cache         *mencache.Mencache
	Repository    *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		BannerQuery: NewBannerQueryService(&BannerQueryServiceDeps{
			Observability:    deps.Observability,
			Cache:            deps.Cache.BannerQueryCache,
			BannerRepository: deps.Repository.BannerQuery,
			Logger:           deps.Logger,
		}),
		BannerCommand: NewBannerCommandService(&BannerCommandServiceDeps{
			Observability:    deps.Observability,
			Cache:            deps.Cache.BannerCommandCache,
			BannerRepository: deps.Repository.BannerCommand,
			Logger:           deps.Logger,
		}),
	}
}
