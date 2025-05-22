package handler

import "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	MerchantPolicy MerchantPoliciesHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		MerchantPolicy: NewMerchantPolicyHandleGrpc(deps.Service),
	}
}
