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

type orderStatsByMerchantService struct {
	observability                  observability.TraceLoggerObservability
	cache                          cache.OrderStatsByMerchantCache
	orderStatsByMerchantRepository repository.OrderStatsByMerchantRepository
	logger                         logger.LoggerInterface
}

type OrderStatsByMerchantServiceDeps struct {
	Observability                  observability.TraceLoggerObservability
	Cache                          cache.OrderStatsByMerchantCache
	OrderStatsByMerchantRepository repository.OrderStatsByMerchantRepository
	Logger                         logger.LoggerInterface
}

func NewOrderStatsByMerchantService(deps *OrderStatsByMerchantServiceDeps) OrderStatsByMerchantService {
	return &orderStatsByMerchantService{
		observability:                  deps.Observability,
		cache:                          deps.Cache,
		orderStatsByMerchantRepository: deps.OrderStatsByMerchantRepository,
		logger:                         deps.Logger,
	}
}

func (s *orderStatsByMerchantService) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error) {
	const method = "FindMonthlyTotalRevenueByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTotalRevenueByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly total revenue by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetMonthlyTotalRevenueByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalRevenueByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetMonthlyTotalRevenueByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched monthly total revenue by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *orderStatsByMerchantService) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, error) {
	const method = "FindYearlyTotalRevenueByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTotalRevenueByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved yearly total revenue by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetYearlyTotalRevenueByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalRevenueByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetYearlyTotalRevenueByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched yearly total revenue by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *orderStatsByMerchantService) FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, error) {
	const method = "FindMonthlyOrderByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly orders by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetMonthlyOrderByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyOrderByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetMonthlyOrderByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched monthly orders by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}

func (s *orderStatsByMerchantService) FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, error) {
	const method = "FindYearlyOrderByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchantID", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved yearly orders by merchant from cache",
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID))
		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetYearlyOrderByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyOrderByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("merchantID", req.MerchantID),
		)
	}

	s.cache.SetYearlyOrderByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched yearly orders by merchant from repository",
		zap.Int("year", req.Year),
		zap.Int("merchantID", req.MerchantID))

	return res, nil
}
