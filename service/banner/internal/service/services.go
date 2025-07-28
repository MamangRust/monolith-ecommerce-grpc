package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	BannerQuery   BannerQueryService
	BannerCommand BannerCommandService
}

type Deps struct {
	Repositories *repository.Repositories
	ErrorHandler *errorhandler.ErroHandler
	Mencache     *mencache.Mencache
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	bannerMapper := response_service.NewBannerResponseMapper()

	return &Service{
		BannerQuery:   NewBannerQueryService(deps.ErrorHandler.BannerQueryError, deps.Mencache.BannerQueryCache, deps.Repositories.BannerQuery, deps.Logger, bannerMapper),
		BannerCommand: NewBannerCommandService(deps.ErrorHandler.BannerCommandError, deps.Mencache.BannerCommandCache, deps.Repositories.BannerCommand, deps.Logger, bannerMapper),
	}
}
