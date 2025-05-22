package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	BannerQuery   BannerQueryService
	BannerCommand BannerCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	bannerMapper := response_service.NewBannerResponseMapper()

	return &Service{
		BannerQuery:   NewBannerQueryService(deps.Ctx, deps.Repositories.BannerQuery, deps.Logger, bannerMapper),
		BannerCommand: NewBannerCommandService(deps.Ctx, deps.Repositories.BannerCommand, deps.Logger, bannerMapper),
	}
}
