package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantQuery           pb.MerchantQueryServiceServer
	MerchantCommandHandler  pb.MerchantCommandServiceServer
	MerchantDocumentQuery   pb.MerchantDocumentQueryServiceServer
	MerchantDocumentCommand pb.MerchantDocumentCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantQuery:           NewMerchantQueryHandler(deps.Service.MerchantQuery, deps.Logger),
		MerchantCommandHandler:  NewMerchantCommandHandler(deps.Service.MerchantCommand, deps.Logger),
		MerchantDocumentQuery:   NewMerchantDocumentQueryHandler(deps.Service.MerchantDocumentQuery, deps.Logger),
		MerchantDocumentCommand: NewMerchantDocumentCommandHandler(deps.Service.MerchantDocumentCommand, deps.Logger),
	}
}
