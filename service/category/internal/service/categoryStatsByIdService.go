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

type categoryStatsByIdService struct {
	ctx                         context.Context
	mencache                    mencache.CategoryStatsByIdCache
	errorhandler                errorhandler.CategoryStatsByIdError
	trace                       trace.Tracer
	categoryStatsByIdRepository repository.CategoryStatsByIdRepository
	logger                      logger.LoggerInterface
	mapping                     response_service.CategoryResponseMapper
	requestCounter              *prometheus.CounterVec
	requestDuration             *prometheus.HistogramVec
}

func NewCategoryStatsByIdService(ctx context.Context,
	mencache mencache.CategoryStatsByIdCache,
	errorhandler errorhandler.CategoryStatsByIdError,
	categoryStatsByIdRepository repository.CategoryStatsByIdRepository, logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryStatsByIdService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "category_stats_by_id_service_request_total",
			Help: "Total number of requests to the CategoryStatsByIdService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "category_stats_by_id_service_request_duration_seconds",
			Help:    "Duration of requests to the CategoryStatsByIdService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &categoryStatsByIdService{
		ctx:                         ctx,
		mencache:                    mencache,
		errorhandler:                errorhandler,
		trace:                       otel.Tracer("category-stats-by-id-service"),
		categoryStatsByIdRepository: categoryStatsByIdRepository,
		logger:                      logger,
		mapping:                     mapping,
		requestCounter:              requestCounter,
		requestDuration:             requestDuration,
	}
}

func (s *categoryStatsByIdService) FindMonthlyTotalPriceById(req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTotalPriceById"

	year := req.Year
	month := req.Month

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthTotalPriceByIdCache(req); found {
		logSuccess("Successfully fetched monthly total price by ID from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetMonthlyTotalPriceById(req)

	if err != nil {
		return s.errorhandler.HandleMonthTotalPriceError(err, method, "FAILED_FIND_MONTHLY_TOTAL_PRICE_BY_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryMonthlyTotalPrices(res)

	s.mencache.SetCachedMonthTotalPriceByIdCache(req, so)

	logSuccess("Successfully fetched monthly total price by ID", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *categoryStatsByIdService) FindYearlyTotalPriceById(req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	const method = "FindYearlyTotalPriceById"

	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearTotalPriceByIdCache(req); found {
		logSuccess("Successfully fetched yearly total price by ID from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetYearlyTotalPricesById(req)

	if err != nil {
		return s.errorhandler.HandleYearTotalPriceError(err, method, "FAILED_FIND_YEARLY_TOTAL_PRICE_BY_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryYearlyTotalPrices(res)

	s.mencache.SetCachedYearTotalPriceByIdCache(req, so)

	logSuccess("Successfully fetched yearly total price by ID", zap.Int("year", year))

	return so, nil
}

func (s *categoryStatsByIdService) FindMonthPriceById(req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	const method = "FindMonthPriceById"

	year := req.Year
	category_id := req.CategoryID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("category.id", category_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthPriceByIdCache(req); found {
		s.logger.Debug("Successfully fetched monthly category prices by ID from cache", zap.Int("year", year), zap.Int("category.id", category_id))
		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetMonthPriceById(req)

	if err != nil {
		return s.errorhandler.HandleMonthPrice(err, method, "FAILED_FIND_MONTH_PRICE_BY_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryMonthlyPrices(res)

	s.mencache.SetCachedMonthPriceByIdCache(req, so)

	logSuccess("Successfully fetched monthly category prices by ID", zap.Int("year", year), zap.Int("category.id", category_id))

	return so, nil
}

func (s *categoryStatsByIdService) FindYearPriceById(req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	const method = "FindYearPriceById"

	year := req.Year
	category_id := req.CategoryID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("category.id", category_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearPriceByIdCache(req); found {
		logSuccess("Successfully fetched yearly category prices by ID from cache", zap.Int("year", year), zap.Int("category.id", category_id))

		return data, nil
	}

	res, err := s.categoryStatsByIdRepository.GetYearPriceById(req)

	if err != nil {
		return s.errorhandler.HandleYearPrice(err, method, "FAILED_FIND_YEAR_PRICE_BY_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCategoryYearlyPrices(res)

	s.mencache.SetCachedYearPriceByIdCache(req, so)

	logSuccess("Successfully fetched yearly category prices by ID", zap.Int("year", year), zap.Int("category.id", category_id))

	return so, nil
}

func (s *categoryStatsByIdService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *categoryStatsByIdService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
