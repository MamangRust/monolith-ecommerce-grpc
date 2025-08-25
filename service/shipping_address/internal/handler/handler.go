package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Shipping ShippingAddressHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Shipping: NewShippingAddressHandleGrpc(deps.Service, deps.Logger),
	}
}
