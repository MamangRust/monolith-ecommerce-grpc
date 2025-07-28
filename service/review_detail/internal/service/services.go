package service

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	ReviewDetailQuery   ReviewDetailQueryService
	ReviewDetailCommand ReviewDetailCommandService
}

type Deps struct {
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewReviewDetailResponseMapper()

	return &Service{
		ReviewDetailQuery:   NewReviewDetailQueryService(deps.Mencache.ReviewDetailQueryCache, deps.ErrorHandler.ReviewDetailQueryError, deps.Repositories.ReviewDetailQuery, mapper, deps.Logger),
		ReviewDetailCommand: NewReviewDetailCommandService(deps.Mencache.ReviewDetailCommandCache, deps.ErrorHandler.ReviewDetailCommandError, deps.Repositories.ReviewDetailQuery, deps.Repositories.ReviewDetailCommand, mapper, deps.Logger),
	}
}
