package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	MerchantDetailQuery       MerchantDetailQueryService
	MerchantDetailCommand     MerchantDetailCommandService
	MerchantSocialLinkCommand MerchantSocialLinkCommandService
}

type Deps struct {
	Cache         *mencache.Mencache
	Repository    *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		MerchantDetailQuery: NewMerchantDetailQueryService(&MerchantDetailQueryServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.MerchantDetailQueryCache,
			Repository:    deps.Repository.MerchantDetailQuery,
			Logger:        deps.Logger,
		}),
		MerchantDetailCommand: NewMerchantDetailCommandService(&MerchantDetailCommandServiceDeps{
			Observability:            deps.Observability,
			Cache:                    deps.Cache.MerchantDetailCommandCache,
			MerchantDetailRepository: deps.Repository.MerchantDetailCommand,
			MerchantQueryRepository:  deps.Repository.MerchantQuery,
			Logger:                   deps.Logger,
		}),
		MerchantSocialLinkCommand: NewMerchantSocialLinkCommandService(&MerchantSocialLinkCommandServiceDeps{
			Observability: deps.Observability,
			Repository:    deps.Repository.MerchantSocialLinkCommand,
			Logger:        deps.Logger,
		}),
	}
}
