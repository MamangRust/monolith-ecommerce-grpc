package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"

	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type categoryStatsService struct {
	observability           observability.TraceLoggerObservability
	cache                   cache.CategoryStatsCache
	categoryStatsRepository repository.CategoryStatsRepository
	logger                  logger.LoggerInterface
}

type CategoryStatsServiceDeps struct {
	Observability           observability.TraceLoggerObservability
	Cache                   cache.CategoryStatsCache
	CategoryStatsRepository repository.CategoryStatsRepository
	Logger                  logger.LoggerInterface
}

func NewCategoryStatsService(
	deps *CategoryStatsServiceDeps) *categoryStatsService {

	return &categoryStatsService{
		cache:                   deps.Cache,
		categoryStatsRepository: deps.CategoryStatsRepository,
		logger:                  deps.Logger,
		observability:           deps.Observability,
	}
}

func (s *categoryStatsService) FindMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, error) {
	const method = "FindMonthlyTotalPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthTotalPriceCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly total price from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.categoryStatsRepository.GetMonthlyTotalPrice(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalPriceRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
		)
	}

	s.cache.SetCachedMonthTotalPriceCache(ctx, req, res)

	logSuccess("Successfully fetched monthly total price", zap.Int("year", req.Year), zap.Int("month", req.Month))
	return res, nil
}

func (s *categoryStatsService) FindYearlyTotalPrice(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, error) {
	const method = "FindYearlyTotalPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearTotalPriceCache(ctx, year); found {
		logSuccess("Successfully retrieved yearly total price from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.categoryStatsRepository.GetYearlyTotalPrices(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalPriceRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearTotalPriceCache(ctx, year, res)

	logSuccess("Successfully fetched yearly total price", zap.Int("year", year))
	return res, nil
}

func (s *categoryStatsService) FindMonthPrice(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, error) {
	const method = "FindMonthPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthPriceCache(ctx, year); found {
		logSuccess("Successfully retrieved month price from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.categoryStatsRepository.GetMonthPrice(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyCategoryRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedMonthPriceCache(ctx, year, res)

	logSuccess("Successfully fetched month price", zap.Int("year", year))
	return res, nil
}

func (s *categoryStatsService) FindYearPrice(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, error) {
	const method = "FindYearPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearPriceCache(ctx, year); found {
		logSuccess("Successfully retrieved year price from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.categoryStatsRepository.GetYearPrice(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyCategoryRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("year", year),
		)
	}

	s.cache.SetCachedYearPriceCache(ctx, year, res)

	logSuccess("Successfully fetched year price", zap.Int("year", year))
	return res, nil
}
