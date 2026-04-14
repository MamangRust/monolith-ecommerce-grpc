package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantService
}

type Deps struct {
	Kafka         *kafka.Kafka
	Cache         *mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		TransactionQuery: NewTransactionQueryService(&TransactionQueryServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.TransactionQueryCache,
			Repository:    deps.Repositories.TransactionQuery,
			Logger:        deps.Logger,
		}),
		TransactionCommand: NewTransactionCommandService(&TransactionCommandServiceDeps{
			Observability:      deps.Observability,
			Kafka:              deps.Kafka,
			Cache:              deps.Cache.TransactionCommandCache,
			TransactionQuery:   deps.Repositories.TransactionQuery,
			TransactionCommand: deps.Repositories.TransactionCommand,
			UserQuery:          deps.Repositories.UserQuery,
			MerchantQuery:      deps.Repositories.MerchantQuery,
			OrderQuery:         deps.Repositories.OrderQuery,
			OrderItem:          deps.Repositories.OrderItem,
			ShippingAddress:    deps.Repositories.ShippingAddress,
			Logger:             deps.Logger,
		}),
		TransactionStats: NewTransactionStatsService(&TransactionStatsServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.TransactionStatsCache,
			Repository:    deps.Repositories.TransactionStats,
			Logger:        deps.Logger,
		}),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantService(&TransactionStatsByMerchantServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.TransactionStatsByMerchantCache,
			Repository:    deps.Repositories.StatsByMerchant,
			Logger:        deps.Logger,
		}),
	}
}
