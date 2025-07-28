package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	SliderQuery   SliderQueryService
	SliderCommand SliderCommandService
}

type Deps struct {
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewSliderResponseMapper()

	return &Service{
		SliderQuery:   NewSliderQueryService(deps.ErrorHandler.SliderQueryError, deps.Mencache.SliderQueryCache, deps.Repositories.SliderQuery, deps.Logger, mapper),
		SliderCommand: NewSliderCommandService(deps.ErrorHandler.SliderCommandError, deps.Repositories.SliderCommand, deps.Logger, mapper),
	}
}
