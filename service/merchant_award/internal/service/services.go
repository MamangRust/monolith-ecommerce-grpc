package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	MerchantAwardQuery   MerchantAwardQueryService
	MerchantAwardCommand MerchantAwardCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	merchantMapper := response_service.NewMerchantAwardResponseMapper()

	return &Service{
		MerchantAwardQuery:   NewMerchantAwardQueryService(deps.Ctx, deps.Repositories.MerchantAwardQuery, deps.Logger, merchantMapper),
		MerchantAwardCommand: NewMerchantAwardCommandService(deps.Ctx, deps.Repositories.MerchantAwardCommand, deps.Logger, merchantMapper),
	}
}
