package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	MerchantAwardQuery   MerchantAwardQueryService
	MerchantAwardCommand MerchantAwardCommandService
}

type Deps struct {
	Cache         *mencache.Mencache
	Repository    *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		MerchantAwardQuery: NewMerchantAwardQueryService(&MerchantAwardQueryServiceDeps{
			Observability:           deps.Observability,
			Cache:                   deps.Cache.MerchantAwardQueryCache,
			MerchantAwardRepository: deps.Repository.MerchantAwardQuery,
			Logger:                  deps.Logger,
		}),
		MerchantAwardCommand: NewMerchantAwardCommandService(&MerchantAwardCommandServiceDeps{
			Observability:           deps.Observability,
			Cache:                   deps.Cache.MerchantAwardCommandCache,
			MerchantAwardRepository: deps.Repository.MerchantAwardCommand,
			Logger:                  deps.Logger,
		}),
	}
}
