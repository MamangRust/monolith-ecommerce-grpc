package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type categoryStatsByMerchantService struct {
	ctx                               context.Context
	trace                             trace.Tracer
	categoryStatsByMerchantRepository repository.CategoryStatsByMerchantRepository
	logger                            logger.LoggerInterface
	mapping                           response_service.CategoryResponseMapper
	requestCounter                    *prometheus.CounterVec
	requestDuration                   *prometheus.HistogramVec
}

func NewCategoryStatsByMerchantService(ctx context.Context, categoryStatsByMerchantRepository repository.CategoryStatsByMerchantRepository, logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryStatsByMerchantService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "category_stats_by_merchant_service_request_total",
			Help: "Total number of requests to the CategoryStatsByMerchantService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "category_stats_by_merchant_service_request_duration_seconds",
			Help:    "Duration of requests to the CategoryStatsByMerchantService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &categoryStatsByMerchantService{
		ctx:                               ctx,
		trace:                             otel.Tracer("category-stats-by-id-service"),
		categoryStatsByMerchantRepository: categoryStatsByMerchantRepository,
		logger:                            logger,
		mapping:                           mapping,
		requestCounter:                    requestCounter,
		requestDuration:                   requestDuration,
	}
}

func (s *categoryStatsByMerchantService) FindMonthlyTotalPriceByMerchant(req *requests.MonthTotalPriceMerchant) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalPriceByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalPriceByMerchant")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("find monthly total price by merchant",
		zap.Int("year", year),
		zap.Int("month", month))

	res, err := s.categoryStatsByMerchantRepository.GetMonthlyTotalPriceByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_PRICE_BY_MERCHANT")

		s.logger.Error("failed to get monthly total sales",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get monthly total sales")

		status = "failed_find_monthly_total_price_by_merchant"

		return nil, category_errors.ErrFailedFindMonthlyTotalPriceByMerchant
	}

	return s.mapping.ToCategoryMonthlyTotalPrices(res), nil
}

func (s *categoryStatsByMerchantService) FindYearlyTotalPriceByMerchant(req *requests.YearTotalPriceMerchant) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalPriceByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalPriceByMerchant")
	defer span.End()

	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("find yearly total price by merchant",
		zap.Int("year", year))

	res, err := s.categoryStatsByMerchantRepository.GetYearlyTotalPricesByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_PRICE_BY_MERCHANT")

		s.logger.Error("failed to get yearly total sales",
			zap.Int("year", year),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly total sales")

		status = "failed_find_yearly_total_price_by_merchant"

		return nil, category_errors.ErrFailedFindYearlyTotalPriceByMerchant
	}

	return s.mapping.ToCategoryYearlyTotalPrices(res), nil
}

func (s *categoryStatsByMerchantService) FindMonthPriceByMerchant(req *requests.MonthPriceMerchant) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthPriceByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthPriceByMerchant")
	defer span.End()

	year := req.Year
	merchant_id := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("find monthly category prices by merchant",
		zap.Int("year", year),
		zap.Int("merchant_id", merchant_id))

	res, err := s.categoryStatsByMerchantRepository.GetMonthPriceByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_PRICE_BY_MERCHANT")

		s.logger.Error("failed to get monthly category prices by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get monthly category prices by merchant")

		status = "failed_find_month_price_by_merchant"

		return nil, category_errors.ErrFailedFindMonthPriceByMerchant
	}

	return s.mapping.ToCategoryMonthlyPrices(res), nil
}

func (s *categoryStatsByMerchantService) FindYearPriceByMerchant(req *requests.YearPriceMerchant) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearPriceByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearPriceByMerchant")
	defer span.End()

	year := req.Year
	merchant_id := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("find yearly category prices by merchant",
		zap.Int("year", year),
		zap.Int("merchant_id", merchant_id))

	res, err := s.categoryStatsByMerchantRepository.GetYearPriceByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_PRICE_BY_MERCHANT")

		s.logger.Error("failed to get yearly category prices by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly category prices by merchant")

		status = "failed_find_year_price_by_merchant"

		return nil, category_errors.ErrFailedFindYearPriceByMerchant
	}

	return s.mapping.ToCategoryYearlyPrices(res), nil
}

func (s *categoryStatsByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
