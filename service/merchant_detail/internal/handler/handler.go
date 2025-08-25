package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantDetail     MerchantDetailHandleGrpc
	MerchantSocialLink MerchantSocialLinkHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	mapper := protomapper.NewMerchantProtoMaper()
	mapperDetail := protomapper.NewMerchantDetailProtoMapper()
	mapperSocial := protomapper.NewMerchantSocialLinkProtoMapper()

	return &Handler{
		MerchantDetail:     NewMerchantDetailHandleGrpc(deps.Service, mapperDetail, mapper, deps.Logger),
		MerchantSocialLink: NewMerchantSocialLinkHandleGrpc(deps.Service, mapperSocial, deps.Logger),
	}
}
