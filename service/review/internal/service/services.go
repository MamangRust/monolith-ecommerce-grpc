package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-review/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-review/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	ReviewQuery   ReviewQueryService
	ReviewCommand ReviewCommandService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewReviewResponseMapper()
	return &Service{
		ReviewQuery:   NewReviewQueryService(deps.Ctx, deps.Mencache.ReviewQueryCache, deps.ErrorHandler.ReviewQueryError, deps.Repositories.ReviewQuery, mapper, deps.Logger),
		ReviewCommand: NewReviewCommandService(deps.Ctx, deps.ErrorHandler.ReviewCommandError, deps.Repositories.ProductQuery, deps.Repositories.UserQuery, deps.Repositories.ReviewQuery, deps.Repositories.ReviewCommand, mapper, deps.Logger),
	}
}
