package handler

import (
	"github.com/MamangRust/monolith-ecommerce-merchant_policy/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	pbmerchant_policy "github.com/MamangRust/monolith-ecommerce-pb/merchant_policy"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantPolicyQuery   pbmerchant_policy.MerchantPolicyQueryServiceServer
	MerchantPolicyCommand pbmerchant_policy.MerchantPolicyCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantPolicyQuery:   NewMerchantPolicyQueryHandler(deps.Service.MerchantPoliciesQuery, deps.Logger),
		MerchantPolicyCommand: NewMerchantPolicyCommandHandler(deps.Service.MerchantPoliciesCommand, deps.Logger),
	}
}
