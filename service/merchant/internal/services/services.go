package services

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"golang.org/x/net/context"
)

type Service struct {
	MerchantQuery           MerchantQueryService
	MerchantCommand         MerchantCommandService
	MerchantDocumentCommand MerchantDocumentCommandService
	MerchantDocumentQuery   MerchantDocumentQueryService
}

type Deps struct {
	Kafka        kafka.Kafka
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	merchantMapper := response_service.NewMerchantResponseMapper()
	merchantDocument := response_service.NewMerchantDocumentResponseMapper()

	return &Service{
		MerchantQuery:           NewMerchantQueryService(deps.Ctx, deps.Repositories.MerchantQuery, deps.Logger, merchantMapper),
		MerchantCommand:         NewMerchantCommandService(deps.Kafka, deps.Ctx, deps.Repositories.UserQuery, deps.Repositories.MerchantQuery, deps.Repositories.MerchantCommand, deps.Logger, merchantMapper),
		MerchantDocumentCommand: NewMerchantDocumentCommandService(deps.Kafka, deps.Ctx, deps.Repositories.MerchantDocumentCommand, deps.Repositories.MerchantQuery, deps.Repositories.UserQuery, deps.Logger, merchantDocument),
		MerchantDocumentQuery:   NewMerchantDocumentQueryService(deps.Ctx, deps.Repositories.MerchantDocumentQuery, deps.Logger, merchantDocument),
	}
}
