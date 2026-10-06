package handler

import (
	"github.com/MamangRust/monolith-ecommerce-merchant/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pbmerchant_document "github.com/MamangRust/monolith-ecommerce-pb/merchant_document"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantQuery           pbmerchant.MerchantQueryServiceServer
	MerchantCommandHandler  pbmerchant.MerchantCommandServiceServer
	MerchantDocumentQuery   pbmerchant_document.MerchantDocumentQueryServiceServer
	MerchantDocumentCommand pbmerchant_document.MerchantDocumentCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantQuery:           NewMerchantQueryHandler(deps.Service.MerchantQuery, deps.Logger),
		MerchantCommandHandler:  NewMerchantCommandHandler(deps.Service.MerchantCommand, deps.Logger),
		MerchantDocumentQuery:   NewMerchantDocumentQueryHandler(deps.Service.MerchantDocumentQuery, deps.Logger),
		MerchantDocumentCommand: NewMerchantDocumentCommandHandler(deps.Service.MerchantDocumentCommand, deps.Logger),
	}
}
