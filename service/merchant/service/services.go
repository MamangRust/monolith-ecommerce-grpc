package service

import (
	merchantCache "github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

type Service struct {
	MerchantQuery           MerchantQueryService
	MerchantCommand         MerchantCommandService
	MerchantDocumentCommand MerchantDocumentCommandService
	MerchantDocumentQuery   MerchantDocumentQueryService
}

type Deps struct {
	Kafka         *kafka.Kafka
	Repositories  *repository.Repositories
	Mencache      *merchantCache.Mencache
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		MerchantQuery: NewMerchantQueryService(&MerchantQueryServiceDeps{
			Cache:              deps.Mencache.MerchantQueryCache,
			MerchantRepository: deps.Repositories.MerchantQuery,
			Logger:             deps.Logger,
			Observability:      deps.Observability,
		}),
		MerchantCommand: NewMerchantCommandService(&MerchantCommandServiceDeps{
			Kafka:              deps.Kafka,
			Cache:              deps.Mencache.MerchantCommandCache,
			MerchantRepository: deps.Repositories.MerchantCommand,
			MerchantQuery:      deps.Repositories.MerchantQuery,
			UserRepository:     deps.Repositories.UserQuery,
			Logger:             deps.Logger,
			Observability:      deps.Observability,
		}),
		MerchantDocumentCommand: NewMerchantDocumentCommandService(&MerchantDocumentCommandServiceDeps{
			Cache:         deps.Mencache.MerchantDocumentCommandCache,
			Repository:    deps.Repositories.MerchantDocumentCommand,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		MerchantDocumentQuery: NewMerchantDocumentQueryService(&MerchantDocumentQueryServiceDeps{
			Cache:         deps.Mencache.MerchantDocumentQueryCache,
			Repository:    deps.Repositories.MerchantDocumentQuery,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
	}
}
