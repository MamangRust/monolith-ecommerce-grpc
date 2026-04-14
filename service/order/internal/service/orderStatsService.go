package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderStatsService struct {
	observability        observability.TraceLoggerObservability
	cache                cache.OrderStatsCache
	orderStatsRepository repository.OrderStatsRepository
	logger               logger.LoggerInterface
}

type OrderStatsServiceDeps struct {
	Observability        observability.TraceLoggerObservability
	Cache                cache.OrderStatsCache
	OrderStatsRepository repository.OrderStatsRepository
	Logger               logger.LoggerInterface
}

func NewOrderStatsService(deps *OrderStatsServiceDeps) OrderStatsService {
	return &orderStatsService{
		observability:        deps.Observability,
		cache:                deps.Cache,
		orderStatsRepository: deps.OrderStatsRepository,
		logger:               deps.Logger,
	}
}

func (s *orderStatsService) FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, error) {
	const method = "FindMonthlyTotalRevenue"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTotalRevenueCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly total revenue from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.orderStatsRepository.GetMonthlyTotalRevenue(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalRevenueRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetMonthlyTotalRevenueCache(ctx, req, res)

	logSuccess("Successfully fetched monthly total revenue from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	return res, nil
}

func (s *orderStatsService) FindYearlyTotalRevenue(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, error) {
	const method = "FindYearlyTotalRevenue"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTotalRevenueCache(ctx, year); found {
		logSuccess("Successfully retrieved yearly total revenue from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderStatsRepository.GetYearlyTotalRevenue(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalRevenueRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyTotalRevenueCache(ctx, year, res)

	logSuccess("Successfully fetched yearly total revenue from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *orderStatsService) FindMonthlyOrder(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, error) {
	const method = "FindMonthlyOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyOrderCache(ctx, year); found {
		logSuccess("Successfully retrieved monthly orders from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderStatsRepository.GetMonthlyOrder(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyOrderRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetMonthlyOrderCache(ctx, year, res)

	logSuccess("Successfully fetched monthly orders from repository",
		zap.Int("year", year))

	return res, nil
}

func (s *orderStatsService) FindYearlyOrder(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, error) {
	const method = "FindYearlyOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyOrderCache(ctx, year); found {
		logSuccess("Successfully retrieved yearly orders from cache",
			zap.Int("year", year))
		return data, nil
	}

	res, err := s.orderStatsRepository.GetYearlyOrder(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyOrderRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", year),
		)
	}

	s.cache.SetYearlyOrderCache(ctx, year, res)

	logSuccess("Successfully fetched yearly orders from repository",
		zap.Int("year", year))

	return res, nil
}
