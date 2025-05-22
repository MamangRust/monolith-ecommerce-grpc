package handler

import "github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	Slider SliderHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Slider: NewSliderHandleGrpc(deps.Service),
	}
}
