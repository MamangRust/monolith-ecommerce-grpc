package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	OrderItemQuery OrderItemQueryService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewOrderItemResponseMapper()

	return &Service{
		OrderItemQuery: NewOrderItemQueryService(deps.Ctx, deps.Repositories.OrderItemQuery, deps.Logger, mapper),
	}
}
