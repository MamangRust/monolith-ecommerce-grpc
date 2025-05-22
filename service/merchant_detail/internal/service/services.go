package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	MerchantDetailQuery   MerchantDetailQueryService
	MerchantDetailCommand MerchantDetailCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewMerchantDetailResponseMapper()
	return &Service{
		MerchantDetailQuery:   NewMerchantDetailQueryService(deps.Ctx, deps.Repositories.MerchantDetailQuery, mapper, deps.Logger),
		MerchantDetailCommand: NewMerchantDetailCommandService(deps.Ctx, deps.Repositories.MerchantDetailQuery, deps.Repositories.MerchantDetailCommand, deps.Repositories.MerchantSocialLinkCommand, mapper, deps.Logger),
	}
}
