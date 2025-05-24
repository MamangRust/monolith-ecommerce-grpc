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

type categoryStatsService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	categoryStatsRepository repository.CategoryStatsRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.CategoryResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCategoryStatsService(ctx context.Context, categoryStatsRepository repository.CategoryStatsRepository, logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryStatsService {
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
		trace:                   otel.Tracer("category-stats-service"),
		categoryStatsRepository: categoryStatsRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *categoryStatsService) FindMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalPrice", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalPrice")
	defer span.End()

	year := req.Year
	month := req.Month

	res, err := s.categoryStatsRepository.GetMonthlyTotalPrice(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_PRICE")

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

		status = "failed_find_monthly_total_price"

		return nil, category_errors.ErrFailedFindMonthlyTotalPrice
	}

	return s.mapping.ToCategoryMonthlyTotalPrices(res), nil
}

func (s *categoryStatsService) FindYearlyTotalPrice(year int) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalPrice", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalPrice")
	defer span.End()

	res, err := s.categoryStatsRepository.GetYearlyTotalPrices(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_PRICE")

		s.logger.Error("failed to get yearly total sales",
			zap.Int("year", year),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly total sales")

		status = "failed_find_yearly_total_price"

		return nil, category_errors.ErrFailedFindYearlyTotalPrice
	}

	return s.mapping.ToCategoryYearlyTotalPrices(res), nil
}

func (s *categoryStatsService) FindMonthPrice(year int) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthPrice", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthPrice")
	defer span.End()

	res, err := s.categoryStatsRepository.GetMonthPrice(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_PRICE")

		s.logger.Error("failed to get monthly category prices",
			zap.Int("year", year),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get monthly category prices")

		status = "failed_find_month_price"

		return nil, category_errors.ErrFailedFindMonthPrice
	}

	return s.mapping.ToCategoryMonthlyPrices(res), nil
}

func (s *categoryStatsService) FindYearPrice(year int) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearPrice", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearPrice")
	defer span.End()

	res, err := s.categoryStatsRepository.GetYearPrice(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_PRICE")

		s.logger.Error("failed to get yearly category prices",
			zap.Int("year", year),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly category prices")

		status = "failed_find_year_price"

		return nil, category_errors.ErrFailedFindYearPrice
	}

	return s.mapping.ToCategoryYearlyPrices(res), nil
}

func (s *categoryStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
