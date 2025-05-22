package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	ReviewDetailQuery   ReviewDetailQueryService
	ReviewDetailCommand ReviewDetailCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewReviewDetailResponseMapper()

	return &Service{
		ReviewDetailQuery:   NewReviewDetailQueryService(deps.Ctx, deps.Repositories.ReviewDetailQuery, mapper, deps.Logger),
		ReviewDetailCommand: NewReviewDetailCommandService(deps.Ctx, deps.Repositories.ReviewDetailQuery, deps.Repositories.ReviewDetailCommand, mapper, deps.Logger),
	}
}
