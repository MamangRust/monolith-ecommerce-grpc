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

type categoryStatsByIdService struct {
	ctx                         context.Context
	trace                       trace.Tracer
	categoryStatsByIdRepository repository.CategoryStatsByIdRepository
	logger                      logger.LoggerInterface
	mapping                     response_service.CategoryResponseMapper
	requestCounter              *prometheus.CounterVec
	requestDuration             *prometheus.HistogramVec
}

func NewCategoryStatsByIdService(ctx context.Context, categoryStatsByIdRepository repository.CategoryStatsByIdRepository, logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryStatsByIdService {
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
		trace:                       otel.Tracer("category-stats-by-id-service"),
		categoryStatsByIdRepository: categoryStatsByIdRepository,
		logger:                      logger,
		mapping:                     mapping,
		requestCounter:              requestCounter,
		requestDuration:             requestDuration,
	}
}

func (s *categoryStatsByIdService) FindMonthlyTotalPriceById(req *requests.MonthTotalPriceCategory) ([]*response.CategoriesMonthlyTotalPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalPriceById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalPriceById")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("find monthly total price by ID",
		zap.Int("year", year),
		zap.Int("month", month))

	res, err := s.categoryStatsByIdRepository.GetMonthlyTotalPriceById(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_PRICE_BY_ID")

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

		status = "failed_find_monthly_total_price_by_id"

		return nil, category_errors.ErrFailedFindMonthlyTotalPriceById
	}

	return s.mapping.ToCategoryMonthlyTotalPrices(res), nil
}

func (s *categoryStatsByIdService) FindYearlyTotalPriceById(req *requests.YearTotalPriceCategory) ([]*response.CategoriesYearlyTotalPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalPriceById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalPriceById")
	defer span.End()

	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("find yearly total price by ID",
		zap.Int("year", year))

	res, err := s.categoryStatsByIdRepository.GetYearlyTotalPricesById(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_PRICE_BY_ID")

		s.logger.Error("failed to get yearly total sales",
			zap.Int("year", year),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly total sales")

		status = "failed_find_yearly_total_price_by_id"

		return nil, category_errors.ErrFailedFindYearlyTotalPriceById
	}

	return s.mapping.ToCategoryYearlyTotalPrices(res), nil
}

func (s *categoryStatsByIdService) FindMonthPriceById(req *requests.MonthPriceId) ([]*response.CategoryMonthPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthPriceById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthPriceById")
	defer span.End()

	year := req.Year
	category_id := req.CategoryID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("category_id", category_id),
	)

	s.logger.Debug("find monthly category prices by ID",
		zap.Int("year", year),
		zap.Int("category_id", category_id))

	res, err := s.categoryStatsByIdRepository.GetMonthPriceById(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_PRICE_BY_ID")

		s.logger.Error("failed to get monthly category prices by ID",
			zap.Int("year", year),
			zap.Int("category_id", category_id),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get monthly category prices by ID")

		status = "failed_find_month_price_by_id"

		return nil, category_errors.ErrFailedFindMonthPriceById
	}

	return s.mapping.ToCategoryMonthlyPrices(res), nil
}

func (s *categoryStatsByIdService) FindYearPriceById(req *requests.YearPriceId) ([]*response.CategoryYearPriceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearPriceById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearPriceById")
	defer span.End()

	year := req.Year
	category_id := req.CategoryID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("category_id", category_id),
	)

	s.logger.Debug("find yearly category prices by ID",
		zap.Int("year", year),
		zap.Int("category_id", category_id))

	res, err := s.categoryStatsByIdRepository.GetYearPriceById(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_PRICE_BY_ID")

		s.logger.Error("failed to get yearly category prices by ID",
			zap.Int("year", year),
			zap.Int("category_id", category_id),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly category prices by ID")

		status = "failed_find_year_price_by_id"

		return nil, category_errors.ErrFailedFindYearPriceById
	}

	return s.mapping.ToCategoryYearlyPrices(res), nil
}

func (s *categoryStatsByIdService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
