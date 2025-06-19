package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/service"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *service.Service
}

type Handler struct {
	MerchantAward MerchantAwardHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	merchantProto := protomapper.NewMerchantProtoMaper()
	merchantAwardProto := protomapper.NewMerchantAwardProtoMapper()

	return &Handler{
		MerchantAward: NewMerchantAwardHandleGrpc(deps.Service, merchantAwardProto, merchantProto),
	}
}
