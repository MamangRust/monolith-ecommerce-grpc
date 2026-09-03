package service

import (
	"github.com/jackc/pgx/v5/pgxpool"

	mencache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantService
	Outbox                     OutboxService
}

type Deps struct {
	Kafka         *kafka.Kafka
	Pool          *pgxpool.Pool
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
			Pool:               deps.Pool,
			Cache:              deps.Cache.TransactionCommandCache,
			TransactionQuery:   deps.Repositories.TransactionQuery,
			TransactionCommand: deps.Repositories.TransactionCommand,
			UserQuery:          deps.Repositories.UserQuery,
			MerchantQuery:      deps.Repositories.MerchantQuery,
			OrderQuery:         deps.Repositories.OrderQuery,
			OrderItem:          deps.Repositories.OrderItem,
			ShippingAddress:    deps.Repositories.ShippingAddress,
			Outbox:             deps.Repositories.Outbox,
			Logger:             deps.Logger,
		}),
		Outbox: NewOutboxService(deps.Repositories.Outbox, deps.Kafka, deps.Logger),
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
