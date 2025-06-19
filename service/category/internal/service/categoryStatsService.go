package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-category/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type categoryStatsService struct {
	ctx                     context.Context
	mencache                mencache.CategoryStatsCache
	errorHandler            errorhandler.CategoryStatsError
	trace                   trace.Tracer
	categoryStatsRepository repository.CategoryStatsRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.CategoryResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCategoryStatsService(ctx context.Context,
	mencache mencache.CategoryStatsCache,
	errorHandler errorhandler.CategoryStatsError,
	categoryStatsRepository repository.CategoryStatsRepository, logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryStatsService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "category_stats_service_request_total",
			Help: "Total number of requests to the CategoryStatsService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "category_stats_service_request_duration_seconds",
			Help:    "Duration of requests to the CategoryStatsService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &categoryStatsService{
		ctx:                     ctx,
		mencache:                mencache,
		errorHandler:            errorHandler,
		trace:                   otel.Tracer("category-stats-service"),
		categoryStatsRepository: categoryStatsRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *categoryStatsService) FindMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTotalPrice"

	year := req.Year
	month := req.Month

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthTotalPriceCache(req); found {
		logSuccess("Successfully fetched monthly total price from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.categoryStatsRepository.GetMonthlyTotalPrice(req)

	if err != nil {
		return s.errorHandler.HandleMonthTotalPriceError(err, method, "FAILED_FIND_MONTHLY_TOTAL_PRICE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryMonthlyTotalPrices(res)

	s.mencache.SetCachedMonthTotalPriceCache(req, so)

	logSuccess("Successfully fetched monthly total price", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *categoryStatsService) FindYearlyTotalPrice(year int) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	const method = "FindYearlyTotalPrice"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearTotalPriceCache(year); found {
		logSuccess("Successfully fetched yearly total price from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.categoryStatsRepository.GetYearlyTotalPrices(year)

	if err != nil {
		return s.errorHandler.HandleYearTotalPriceError(err, method, "FAILED_FIND_YEARLY_TOTAL_PRICE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryYearlyTotalPrices(res)

	s.mencache.SetCachedYearTotalPriceCache(year, so)

	logSuccess("Successfully fetched yearly total price", zap.Int("year", year))

	return so, nil
}

func (s *categoryStatsService) FindMonthPrice(year int) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	const method = "FindMonthPrice"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthPriceCache(year); found {
		logSuccess("Successfully fetched monthly category prices from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.categoryStatsRepository.GetMonthPrice(year)

	if err != nil {
		return s.errorHandler.HandleMonthPrice(err, method, "FAILED_FIND_MONTH_PRICE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryMonthlyPrices(res)

	s.mencache.SetCachedMonthPriceCache(year, so)

	logSuccess("Successfully fetched monthly category prices", zap.Int("year", year))

	return so, nil
}

func (s *categoryStatsService) FindYearPrice(year int) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	const method = "FindYearPrice"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearPriceCache(year); found {
		logSuccess("Successfully fetched yearly category prices from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.categoryStatsRepository.GetYearPrice(year)

	if err != nil {
		return s.errorHandler.HandleYearPrice(err, method, "FAILED_FIND_YEAR_PRICE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryYearlyPrices(res)

	s.mencache.SetCachedYearPriceCache(year, so)

	logSuccess("Successfully fetched yearly category prices", zap.Int("year", year))

	return so, nil
}

func (s *categoryStatsService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return span, end, status, logSuccess
}

func (s *categoryStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
