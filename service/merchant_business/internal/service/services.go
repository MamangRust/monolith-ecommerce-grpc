package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	MerchantBusinessCommand MerchantBusinessCommandService
	MerchantBusinessQuery   MerchantBusinessQueryService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewMerchantBusinessResponseMapper()

	return &Service{
		MerchantBusinessCommand: NewMerchantBusinessCommandService(deps.Ctx, deps.Repositories.MerchantQuery, deps.Repositories.MerchantBusinessCmd, deps.Logger, mapper),
		MerchantBusinessQuery:   NewMerchantBusinessQueryService(deps.Ctx, deps.Repositories.MerchantBusinessQuery, deps.Logger, mapper),
	}
}
