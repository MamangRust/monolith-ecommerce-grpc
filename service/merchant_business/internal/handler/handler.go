package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantBusiness MerchantBusinessHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	mapper := protomapper.NewMerchantProtoMaper()
	mapperBusiness := protomapper.NewMerchantBusinessProtoMapper()

	return &Handler{
		MerchantBusiness: NewMerchantBusinessHandleGrpc(deps.Logger, deps.Service, mapperBusiness, mapper),
	}
}
