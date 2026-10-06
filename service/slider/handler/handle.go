package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-slider/service"

	pbslider "github.com/MamangRust/monolith-ecommerce-pb/slider"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	SliderQuery   pbslider.SliderQueryServiceServer
	SliderCommand pbslider.SliderCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		SliderQuery:   NewSliderQueryHandler(deps.Service.SliderQuery, deps.Logger),
		SliderCommand: NewSliderCommandHandler(deps.Service.SliderCommand, deps.Logger),
	}
}
