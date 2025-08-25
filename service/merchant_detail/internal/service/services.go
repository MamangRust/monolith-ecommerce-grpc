package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	MerchantDetailQuery   MerchantDetailQueryService
	MerchantDetailCommand MerchantDetailCommandService
	MerchantSocialLink    MerchantSocialLinkService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewMerchantDetailResponseMapper()
	mapper_link := response_service.NewMerchantSocialLinkResponseMapper()
	return &Service{
		MerchantDetailQuery:   NewMerchantDetailQueryService(deps.Ctx, deps.ErrorHandler.MerchantDetailQueryError, deps.Mencache.MerchantDetailQueryCache, deps.Repositories.MerchantDetailQuery, mapper, deps.Logger),
		MerchantDetailCommand: NewMerchantDetailCommandService(deps.Ctx, deps.ErrorHandler.MerchantDetailCommandError, deps.ErrorHandler.FileError, deps.Mencache.MerchantDetailCommandCache, deps.Repositories.MerchantDetailQuery, deps.Repositories.MerchantDetailCommand, deps.Repositories.MerchantSocialLinkCommand, mapper, deps.Logger),
		MerchantSocialLink:    NewMerchantSocialLinkService(deps.Ctx, deps.ErrorHandler.MerchantSocialLinkCommandError, deps.Repositories.MerchantSocialLinkCommand, mapper_link, deps.Logger),
	}
}
