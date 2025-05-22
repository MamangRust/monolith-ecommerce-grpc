package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	MerchantPolicyQuery MerchantPoliciesQueryService
	MerchantPolicyCmd   MerchantPoliciesCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewServices(deps Deps) *Service {
	mapper := response_service.NewMerchantPolicyResponseMapper()

	return &Service{
		MerchantPolicyQuery: NewMerchantPolicyQueryService(deps.Ctx, deps.Logger, deps.Repositories.MerchantPolicyQuery, mapper),
		MerchantPolicyCmd:   NewMerchantPolicyCommandService(deps.Ctx, deps.Logger, deps.Repositories.MerchantPolicyCmd, deps.Repositories.MerchantQuery, mapper),
	}
}
