package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/services"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
)

type Deps struct {
	Service services.Service
}

type Handler struct {
	Merchant         MerchantHandleGrpc
	MerchantDocument MerchantDocumentHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	merchantProto := protomapper.NewMerchantProtoMaper()
	merchantDocumentProto := protomapper.NewMerchantDocumentProtoMapper()

	return &Handler{
		Merchant:         NewMerchantHandleGrpc(deps.Service, merchantProto),
		MerchantDocument: NewMerchantDocumentHandleGrpc(deps.Service, merchantDocumentProto),
	}
}
