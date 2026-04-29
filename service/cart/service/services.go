package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-cart/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	CartQuery   CartQueryService
	CartCommand CartCommandService
}

type Deps struct {
	Cache         mencache.CartMencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		CartQuery: NewCartQueryService(&CartQueryServiceDeps{
			Observability:       deps.Observability,
			Mencache:            deps.Cache,
			CartQueryRepository: deps.Repositories.CartQuery,
			Logger:              deps.Logger,
		}),
		CartCommand: NewCartCommandService(&CartCommandServiceDeps{
			Observability:          deps.Observability,
			CartCommandRepository:  deps.Repositories.CartCommand,
			ProductQueryRepository: deps.Repositories.ProductQuery,
			UserQueryRepository:    deps.Repositories.UserQuery,
			Logger:                 deps.Logger,
		}),
	}
}
