package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantPolicyQuery   pb.MerchantPolicyQueryServiceServer
	MerchantPolicyCommand pb.MerchantPolicyCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantPolicyQuery:   NewMerchantPolicyQueryHandler(deps.Service.MerchantPoliciesQuery, deps.Logger),
		MerchantPolicyCommand: NewMerchantPolicyCommandHandler(deps.Service.MerchantPoliciesCommand, deps.Logger),
	}
}
