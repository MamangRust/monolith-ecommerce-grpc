package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/services"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantService
}

type Deps struct {
	Kafka        kafka.Kafka
	Ctx          context.Context
	Repositories *repository.Repositoris
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewTransactionResponseMapper()

	return &Service{
		TransactionQuery:           NewTransactionQueryService(deps.Ctx, deps.Repositories.TransactionQueryRepository, mapper, deps.Logger),
		TransactionCommand:         NewTransactionCommandService(deps.Ctx, deps.Kafka, deps.Repositories.UserQuery, deps.Repositories.MerchantRepository, deps.Repositories.TransactionQueryRepository, deps.Repositories.TransactionCommandRepository, deps.Repositories.OrderQueryRepository, deps.Repositories.OrderItemRepository, deps.Repositories.ShippingAddressQueryRepository, mapper, deps.Logger),
		TransactionStats:           NewTransactionStatsService(deps.Ctx, deps.Repositories.TransactionStatsRepository, mapper, deps.Logger),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantService(deps.Ctx, deps.Repositories.TransactionStatsByMerchant, mapper, deps.Logger),
	}
}
