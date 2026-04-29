package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	MerchantBusinessQuery   MerchantBusinessQueryService
	MerchantBusinessCommand MerchantBusinessCommandService
}

type Deps struct {
	Cache         *mencache.Mencache
	Repository    *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		MerchantBusinessQuery: NewMerchantBusinessQueryService(&MerchantBusinessQueryServiceDeps{
			Observability:              deps.Observability,
			Cache:                      deps.Cache.MerchantBusinessQueryCache,
			MerchantBusinessRepository: deps.Repository.MerchantBusinessQuery,
			Logger:                     deps.Logger,
		}),
		MerchantBusinessCommand: NewMerchantBusinessCommandService(&MerchantBusinessCommandServiceDeps{
			Observability:              deps.Observability,
			Cache:                      deps.Cache.MerchantBusinessCommandCache,
			MerchantBusinessRepository: deps.Repository.MerchantBusinessCommand,
			Logger:                     deps.Logger,
		}),
	}
}
