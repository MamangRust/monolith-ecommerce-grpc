package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Slider SliderHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Slider: NewSliderHandleGrpc(deps.Service, deps.Logger),
	}
}
