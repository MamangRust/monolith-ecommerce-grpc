package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	ShippingAddressQuery   ShippingAddressQueryService
	ShippingAddressCommand ShippingAddressCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewShippingAddressResponseMapper()

	return &Service{
		ShippingAddressQuery:   NewShippingAddressQueryService(deps.Ctx, deps.Repositories.ShippingAddressQuery, deps.Logger, mapper),
		ShippingAddressCommand: NewShippingAddressCommandService(deps.Ctx, deps.Repositories.ShippingAddressCommand, deps.Logger, mapper),
	}
}
