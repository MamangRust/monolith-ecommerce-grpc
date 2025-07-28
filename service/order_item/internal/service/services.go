package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	OrderItemQuery OrderItemQueryService
}

type Deps struct {
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewOrderItemResponseMapper()

	return &Service{
		OrderItemQuery: NewOrderItemQueryService(deps.ErrorHandler.OrderItemQueryError, deps.Mencache.OrderItemQueryCache, deps.Repositories.OrderItemQuery, deps.Logger, mapper),
	}
}
