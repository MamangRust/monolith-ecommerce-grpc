package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	CartQuery   CartQueryService
	CartCommand CartCommandService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencached    *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewCartResponseMapper()

	return &Service{
		CartQuery:   NewCartQueryService(deps.Ctx, deps.ErrorHandler.CartQueryError, deps.Mencached.CartQueryCache, deps.Repositories.CartQuery, deps.Logger, mapper),
		CartCommand: NewCardCommandService(deps.Ctx, deps.ErrorHandler.CartCommandError, deps.Repositories.CartCommand, deps.Repositories.ProductQuery, deps.Repositories.UserQuery, deps.Logger, mapper),
	}
}
