package handler

import "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	Shipping ShippingAddressHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Shipping: NewShippingAddressHandleGrpc(deps.Service),
	}
}
