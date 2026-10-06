package service

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	mencache "github.com/MamangRust/monolith-ecommerce-slider/cache"
	"github.com/MamangRust/monolith-ecommerce-slider/repository"
)

type Service struct {
	SliderQuery   SliderQueryService
	SliderCommand SliderCommandService
}

type Deps struct {
	Repositories  *repository.Repositories
	Mencache      mencache.SliderMencache
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		SliderQuery: NewSliderQueryService(&SliderQueryServiceDeps{
			Repositories:  deps.Repositories.SliderQuery,
			Cache:         deps.Mencache,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		SliderCommand: NewSliderCommandService(&SliderCommandServiceDeps{
			Repositories:  deps.Repositories.SliderCommand,
			Cache:         deps.Mencache,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
	}
}
