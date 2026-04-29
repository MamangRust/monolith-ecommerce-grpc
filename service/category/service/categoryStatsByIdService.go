package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type categoryStatsByIdService struct {
	observability               observability.TraceLoggerObservability
	cache                       cache.CategoryStatsByIdCache
	categoryStatsByIdRepository repository.CategoryStatsByIdRepository
	logger                      logger.LoggerInterface
}

type CategoryStatsByIdServiceDeps struct {
	Observability               observability.TraceLoggerObservability
	Cache                       cache.CategoryStatsByIdCache
	CategoryStatsByIdRepository repository.CategoryStatsByIdRepository
	Logger                      logger.LoggerInterface
}

func NewCategoryStatsByIdService(deps *CategoryStatsByIdServiceDeps) CategoryStatsByIdService {
	return &categoryStatsByIdService{
		observability:               deps.Observability,
		cache:                       deps.Cache,
		categoryStatsByIdRepository: deps.CategoryStatsByIdRepository,
		logger:                      deps.Logger,
	}
}

func (s *categoryStatsByIdService) FindMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*db.GetMonthlyTotalPriceByIdRow, error) {
	const method = "FindMonthlyTotalPriceById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("categoryID", req.CategoryID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthTotalPriceByIdCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly total price by ID from cache",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("categoryID", req.CategoryID))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetMonthlyTotalPriceById(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalPriceByIdRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("categoryID", req.CategoryID),
		)
	}

	s.cache.SetCachedMonthTotalPriceByIdCache(ctx, req, res)

	logSuccess("Successfully fetched monthly total price by ID from repository",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("categoryID", req.CategoryID))

	return res, nil
}

func (s *categoryStatsByIdService) FindYearlyTotalPriceById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*db.GetYearlyTotalPriceByIdRow, error) {
	const method = "FindYearlyTotalPriceById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("categoryID", req.CategoryID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearTotalPriceByIdCache(ctx, req); found {
		logSuccess("Successfully retrieved yearly total price by ID from cache",
			zap.Int("year", req.Year),
			zap.Int("categoryID", req.CategoryID))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetYearlyTotalPricesById(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalPriceByIdRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("categoryID", req.CategoryID),
		)
	}

	s.cache.SetCachedYearTotalPriceByIdCache(ctx, req, res)

	logSuccess("Successfully fetched yearly total price by ID from repository",
		zap.Int("year", req.Year),
		zap.Int("categoryID", req.CategoryID))

	return res, nil
}

func (s *categoryStatsByIdService) FindMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*db.GetMonthlyCategoryByIdRow, error) {
	const method = "FindMonthPriceById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("categoryID", req.CategoryID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthPriceByIdCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly category prices by ID from cache",
			zap.Int("year", req.Year),
			zap.Int("categoryID", req.CategoryID))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetMonthPriceById(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyCategoryByIdRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("categoryID", req.CategoryID),
		)
	}

	s.cache.SetCachedMonthPriceByIdCache(ctx, req, res)

	logSuccess("Successfully fetched monthly category prices by ID from repository",
		zap.Int("year", req.Year),
		zap.Int("categoryID", req.CategoryID))

	return res, nil
}

func (s *categoryStatsByIdService) FindYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*db.GetYearlyCategoryByIdRow, error) {
	const method = "FindYearPriceById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("categoryID", req.CategoryID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearPriceByIdCache(ctx, req); found {
		logSuccess("Successfully retrieved yearly category prices by ID from cache",
			zap.Int("year", req.Year),
			zap.Int("categoryID", req.CategoryID))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetYearPriceById(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyCategoryByIdRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("year", req.Year),
			zap.Int("categoryID", req.CategoryID),
		)
	}

	s.cache.SetCachedYearPriceByIdCache(ctx, req, res)

	logSuccess("Successfully fetched yearly category prices by ID from repository",
		zap.Int("year", req.Year),
		zap.Int("categoryID", req.CategoryID))

	return res, nil
}
