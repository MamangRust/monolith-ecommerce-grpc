package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-ecommerce-grpc-cart/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cartQueryService struct {
	observability       observability.TraceLoggerObservability
	mencache            mencache.CartQueryCache
	cartQueryRepository repository.CartQueryRepository
	logger              logger.LoggerInterface
}

type CartQueryServiceDeps struct {
	Observability       observability.TraceLoggerObservability
	Mencache            mencache.CartQueryCache
	CartQueryRepository repository.CartQueryRepository
	Logger              logger.LoggerInterface
}

func NewCartQueryService(
	deps *CartQueryServiceDeps) *cartQueryService {

	return &cartQueryService{
		mencache:            deps.Mencache,
		cartQueryRepository: deps.CartQueryRepository,
		logger:              deps.Logger,
		observability:       deps.Observability,
	}
}

func (s *cartQueryService) FindAll(ctx context.Context, req *requests.FindAllCarts) ([]*db.GetCartsRow, *int, error) {
	const method = "FindAllCarts"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCartsCache(ctx, req); found {
		logSuccess("Successfully retrieved all cart records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	carts, err := s.cartQueryRepository.FindCarts(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetCartsRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(carts) > 0 {
		totalCount = int(carts[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.mencache.SetCartsCache(ctx, req, carts, &totalCount)

	logSuccess("Successfully fetched all carts",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return carts, &totalCount, nil
}
