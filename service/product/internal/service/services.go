package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	ProductQuery   ProductQueryService
	ProductCommand ProductCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewProductResponseMapper()

	return &Service{
		ProductQuery:   NewProductQueryService(deps.Ctx, deps.Repositories.ProductQuery, mapper, deps.Logger),
		ProductCommand: NewProductCommandService(deps.Ctx, deps.Repositories.CategoryQuery, deps.Repositories.MerchantQuery, deps.Repositories.ProductQuery, deps.Repositories.ProductCommand, mapper, deps.Logger),
	}
}
