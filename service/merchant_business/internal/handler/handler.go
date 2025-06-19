package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/service"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *service.Service
}

type Handler struct {
	MerchantBusiness MerchantBusinessHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	mapper := protomapper.NewMerchantProtoMaper()
	mapperBusiness := protomapper.NewMerchantBusinessProtoMapper()

	return &Handler{
		MerchantBusiness: NewMerchantBusinessHandleGrpc(deps.Service, mapperBusiness, mapper),
	}
}
