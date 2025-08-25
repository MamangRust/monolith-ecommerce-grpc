package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/services"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service *services.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Merchant         MerchantHandleGrpc
	MerchantDocument MerchantDocumentHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	merchantProto := protomapper.NewMerchantProtoMaper()
	merchantDocumentProto := protomapper.NewMerchantDocumentProtoMapper()

	return &Handler{
		Merchant:         NewMerchantHandleGrpc(deps.Service, merchantProto, deps.Logger),
		MerchantDocument: NewMerchantDocumentHandleGrpc(deps.Service, merchantDocumentProto, deps.Logger),
	}
}
