package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantAward MerchantAwardHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	merchantProto := protomapper.NewMerchantProtoMaper()
	merchantAwardProto := protomapper.NewMerchantAwardProtoMapper()

	return &Handler{
		MerchantAward: NewMerchantAwardHandleGrpc(deps.Service, merchantAwardProto, merchantProto, deps.Logger),
	}
}
