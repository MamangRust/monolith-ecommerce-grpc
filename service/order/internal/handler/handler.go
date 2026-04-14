package handler

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	OrderQuery            OrderQueryHandler
	OrderCommand          OrderCommandHandler
	OrderStats            OrderStatsHandler
	OrderStatsByMerchant  OrderStatsByMerchantHandler
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderQuery:            NewOrderQueryHandler(deps.Service.OrderQuery, deps.Logger),
		OrderCommand:          NewOrderCommandHandler(deps.Service.OrderCommand, deps.Logger),
		OrderStats:            NewOrderStatsHandler(deps.Service.OrderStats, deps.Logger),
		OrderStatsByMerchant:  NewOrderStatsByMerchantHandler(deps.Service.OrderStatsByMerchant, deps.Logger),
	}
}
