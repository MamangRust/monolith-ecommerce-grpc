package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/service"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *service.Service
}

type Handler struct {
	MerchantDetail MerchantDetailHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	mapper := protomapper.NewMerchantProtoMaper()
	mapperDetail := protomapper.NewMerchantDetailProtoMapper()

	return &Handler{
		MerchantDetail: NewMerchantDetailHandleGrpc(deps.Service, mapperDetail, mapper),
	}
}
