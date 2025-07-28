package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantService
}

type Deps struct {
	Kafka        *kafka.Kafka
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositoris
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewTransactionResponseMapper()

	return &Service{
		TransactionQuery:           NewTransactionQueryService(deps.Mencache.TransactionQueryCache, deps.ErrorHandler.TransactionQueryError, deps.Repositories.TransactionQueryRepository, mapper, deps.Logger),
		TransactionCommand:         NewTransactionCommandService(deps.Mencache.TransactionCommandCache, deps.ErrorHandler.TransactionCommandError, deps.Kafka, deps.Repositories.UserQuery, deps.Repositories.MerchantRepository, deps.Repositories.TransactionQueryRepository, deps.Repositories.TransactionCommandRepository, deps.Repositories.OrderQueryRepository, deps.Repositories.OrderItemRepository, deps.Repositories.ShippingAddressQueryRepository, mapper, deps.Logger),
		TransactionStats:           NewTransactionStatsService(deps.ErrorHandler.TransactionStatsError, deps.Mencache.TransactionStatsCache, deps.Repositories.TransactionStatsRepository, mapper, deps.Logger),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantService(deps.ErrorHandler.TransactonStatsByMerchantError, deps.Mencache.TransactionStatsByMerchant, deps.Repositories.TransactionStatsByMerchant, mapper, deps.Logger),
	}
}
