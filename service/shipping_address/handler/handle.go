package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shipping-address/service"

	pbshipping_address "github.com/MamangRust/monolith-ecommerce-pb/shipping_address"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	ShippingQuery   pbshipping_address.ShippingQueryServiceServer
	ShippingCommand pbshipping_address.ShippingCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		ShippingQuery:   NewShippingQueryHandler(deps.Service.ShippingAddressQuery, deps.Logger),
		ShippingCommand: NewShippingCommandHandler(deps.Service.ShippingAddressCommand, deps.Logger),
	}
}
