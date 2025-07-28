package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	MerchantAwardQuery   MerchantAwardQueryService
	MerchantAwardCommand MerchantAwardCommandService
}

type Deps struct {
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	merchantMapper := response_service.NewMerchantAwardResponseMapper()

	return &Service{
		MerchantAwardQuery:   NewMerchantAwardQueryService(deps.ErrorHandler.MerchantAwardQueryError, deps.Mencache.MerchantAwardQueryCache, deps.Repositories.MerchantAwardQuery, deps.Logger, merchantMapper),
		MerchantAwardCommand: NewMerchantAwardCommandService(deps.ErrorHandler.MerchantAwardCommandError, deps.Mencache.MerchantAwardCommandCache, deps.Repositories.MerchantAwardCommand, deps.Logger, merchantMapper),
	}
}
