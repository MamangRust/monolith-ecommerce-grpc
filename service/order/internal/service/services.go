package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-order/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	OrderQuery           OrderQueryService
	OrderCommand         OrderCommandService
	OrderStats           OrderStatsService
	OrderStatsByMerchant OrderStatsByMerchantService
}

type Deps struct {
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewOrderResponseMapper()
	return &Service{
		OrderQuery:           NewOrderQueryService(deps.ErrorHandler.OrderQueryError, deps.Mencache.OrderQueryCache, deps.Repositories.OrderQuery, deps.Logger, mapper),
		OrderCommand:         NewOrderCommandService(deps.ErrorHandler.OrderCommandError, deps.Mencache.OrderCommandCache, deps.Repositories.UserQuery, deps.Repositories.ProductQuery, deps.Repositories.ProductCommand, deps.Repositories.OrderQuery, deps.Repositories.OrderCommand, deps.Repositories.OrderItemQuery, deps.Repositories.OrderItemCommand, deps.Repositories.MerchantQuery, deps.Repositories.ShippingAddress, deps.Logger, mapper),
		OrderStats:           NewOrderStatsService(deps.ErrorHandler.OrderStats, deps.Mencache.OrderStatsCache, deps.Repositories.OrderStats, deps.Logger, mapper),
		OrderStatsByMerchant: NewOrderStatsByMerchantService(deps.Mencache.OrderStatsByMerchantCache, deps.ErrorHandler.OrderStatsByMerchant, deps.Repositories.OrderStatsByMerchant, deps.Logger, mapper),
	}
}
