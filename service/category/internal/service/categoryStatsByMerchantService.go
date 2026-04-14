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

type categoryStatsByMerchantService struct {
	observability                     observability.TraceLoggerObservability
	cache                             cache.CategoryStatsByMerchantCache
	categoryStatsByMerchantRepository repository.CategoryStatsByMerchantRepository
	logger                            logger.LoggerInterface
}

type CategoryStatsByMerchantServiceDeps struct {
	Observability                     observability.TraceLoggerObservability
	Cache                             cache.CategoryStatsByMerchantCache
	CategoryStatsByMerchantRepository repository.CategoryStatsByMerchantRepository
	Logger                            logger.LoggerInterface
}

func NewCategoryStatsByMerchantService(
	deps *CategoryStatsByMerchantServiceDeps) *categoryStatsByMerchantService {
	return &categoryStatsByMerchantService{
		cache:                             deps.Cache,
		categoryStatsByMerchantRepository: deps.CategoryStatsByMerchantRepository,
		logger:                            deps.Logger,
		observability:                     deps.Observability,
	}
}

func (s *categoryStatsByMerchantService) FindMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*db.GetMonthlyTotalPriceByMerchantRow, error) {
	const method = "FindMonthlyTotalPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthTotalPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved monthly total price by merchant from cache", zap.Int("merchantID", req.MerchantID), zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetMonthlyTotalPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTotalPriceByMerchantRow](
			s.logger,
			err,
			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
		)
	}

	s.cache.SetCachedMonthTotalPriceByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched monthly total price by merchant", zap.Int("merchantID", req.MerchantID))
	return res, nil
}

func (s *categoryStatsByMerchantService) FindYearlyTotalPriceByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*db.GetYearlyTotalPriceByMerchantRow, error) {
	const method = "FindYearlyTotalPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", req.MerchantID),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearTotalPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved yearly total price by merchant from cache", zap.Int("merchantID", req.MerchantID), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetYearlyTotalPricesByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTotalPriceByMerchantRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
		)
	}

	s.cache.SetCachedYearTotalPriceByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched yearly total price by merchant", zap.Int("merchantID", req.MerchantID))
	return res, nil
}

func (s *categoryStatsByMerchantService) FindMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*db.GetMonthlyCategoryByMerchantRow, error) {
	const method = "FindMonthPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", req.MerchantID),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved month price by merchant from cache", zap.Int("merchantID", req.MerchantID), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetMonthPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyCategoryByMerchantRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
		)
	}

	s.cache.SetCachedMonthPriceByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched month price by merchant", zap.Int("merchantID", req.MerchantID))
	return res, nil
}

func (s *categoryStatsByMerchantService) FindYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*db.GetYearlyCategoryByMerchantRow, error) {
	const method = "FindYearPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", req.MerchantID),
		attribute.Int("year", req.Year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved year price by merchant from cache", zap.Int("merchantID", req.MerchantID), zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetYearPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyCategoryByMerchantRow](
			s.logger,
			err,

			method,
			span,

			zap.Int("merchant_id", req.MerchantID),
		)
	}

	s.cache.SetCachedYearPriceByMerchantCache(ctx, req, res)

	logSuccess("Successfully fetched year price by merchant", zap.Int("merchantID", req.MerchantID))
	return res, nil
}
