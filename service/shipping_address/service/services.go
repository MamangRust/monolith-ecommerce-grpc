package service

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	mencache "github.com/MamangRust/monolith-ecommerce-shipping-address/cache"
	"github.com/MamangRust/monolith-ecommerce-shipping-address/repository"
)

type Service struct {
	ShippingAddressQuery   ShippingAddressQueryService
	ShippingAddressCommand ShippingAddressCommandService
}

type Deps struct {
	Mencache      mencache.ShippingAddressMencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		ShippingAddressQuery: NewShippingAddressQueryService(&ShippingAddressQueryServiceDeps{
			Observability:             deps.Observability,
			Cache:                     deps.Mencache,
			ShippingAddressRepository: deps.Repositories.ShippingAddressQuery,
			Logger:                    deps.Logger,
		}),
		ShippingAddressCommand: NewShippingAddressCommandService(&ShippingAddressCommandServiceDeps{
			Observability:             deps.Observability,
			Cache:                     deps.Mencache,
			ShippingAddressRepository: deps.Repositories.ShippingAddressCommand,
			Logger:                    deps.Logger,
		}),
	}
}
