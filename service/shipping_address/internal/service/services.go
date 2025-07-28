package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	ShippingAddressQuery   ShippingAddressQueryService
	ShippingAddressCommand ShippingAddressCommandService
}

type Deps struct {
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewShippingAddressResponseMapper()

	return &Service{
		ShippingAddressQuery:   NewShippingAddressQueryService(deps.Mencache.ShippingAddressQueryCache, deps.ErrorHandler.ShippingAddressQueryError, deps.Repositories.ShippingAddressQuery, deps.Logger, mapper),
		ShippingAddressCommand: NewShippingAddressCommandService(deps.Mencache.ShippingAddressCommandCache, deps.ErrorHandler.ShippingAddressCommandError, deps.Repositories.ShippingAddressCommand, deps.Logger, mapper),
	}
}
