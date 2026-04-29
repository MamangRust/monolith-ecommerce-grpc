package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	MerchantBusinessQuery   MerchantBusinessQueryHandler
	MerchantBusinessCommand MerchantBusinessCommandHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		MerchantBusinessQuery:   NewMerchantBusinessQueryHandler(deps.Service.MerchantBusinessQuery, deps.Logger),
		MerchantBusinessCommand: NewMerchantBusinessCommandHandler(deps.Service.MerchantBusinessCommand, deps.Logger),
	}
}
