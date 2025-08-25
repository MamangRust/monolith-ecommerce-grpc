package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantPolicy MerchantPoliciesHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantPolicy: NewMerchantPolicyHandleGrpc(deps.Service, deps.Logger),
	}
}
