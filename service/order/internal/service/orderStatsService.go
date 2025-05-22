package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderStatsService struct {
	ctx                  context.Context
	trace                trace.Tracer
	orderStatsRepository repository.OrderStatsRepository
	logger               logger.LoggerInterface
	mapping              response_service.OrderResponseMapper
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewOrderStatsService(
	ctx context.Context,
	orderStatsRepository repository.OrderStatsRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderResponseMapper,
) *orderStatsService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_stats_service_request_count",
			Help: "Total number of requests to the OrderStatsService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_stats_service_request_duration",
			Help:    "Histogram of request durations for the OrderStatsService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &orderStatsService{
		ctx:                  ctx,
		trace:                otel.Tracer("order-stats-service"),
		orderStatsRepository: orderStatsRepository,
		logger:               logger,
		mapping:              mapping,
		requestCounter:       requestCounter,
		requestDuration:      requestDuration,
	}
}

func (s *orderStatsService) FindMonthlyTotalRevenue(req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalRevenue", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalRevenue")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("find monthly total revenue",
		zap.Int("year", year),
		zap.Int("month", month))

	res, err := s.orderStatsRepository.GetMonthlyTotalRevenue(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_REVENUE")

		s.logger.Error("failed to get monthly total revenue",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to get monthly total revenue")

		status = "failed_find_monthly_total_revenue"

		return nil, order_errors.ErrFailedFindMonthlyTotalRevenue
	}

	return s.mapping.ToOrderMonthlyTotalRevenues(res), nil
}

func (s *orderStatsService) FindYearlyTotalRevenue(year int) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalRevenue", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalRevenue")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	res, err := s.orderStatsRepository.GetYearlyTotalRevenue(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_REVENUE")

		s.logger.Error("failed to get yearly total revenue",
			zap.Int("year", year),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to get yearly total revenue")

		status = "failed_find_yearly_total_revenue"

		return nil, order_errors.ErrFailedFindYearlyTotalRevenue
	}

	return s.mapping.ToOrderYearlyTotalRevenues(res), nil
}

func (s *orderStatsService) FindMonthlyOrder(year int) ([]*response.OrderMonthlyResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyOrder", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyOrder")
	defer span.End()

	res, err := s.orderStatsRepository.GetMonthlyOrder(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_ORDER")

		s.logger.Error("failed to get monthly orders",
			zap.Int("year", year),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to get monthly orders")

		status = "failed_find_monthly_order"

		return nil, order_errors.ErrFailedFindMonthlyOrder
	}

	return s.mapping.ToOrderMonthlyPrices(res), nil
}

func (s *orderStatsService) FindYearlyOrder(year int) ([]*response.OrderYearlyResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyOrder", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	res, err := s.orderStatsRepository.GetYearlyOrder(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_ORDER")

		s.logger.Error("failed to get yearly orders",
			zap.Int("year", year),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to get yearly orders")

		span.SetStatus(codes.Error, "Failed to get yearly orders")

		status = "failed_find_yearly_order"

		return nil, order_errors.ErrFailedFindYearlyOrder
	}

	return s.mapping.ToOrderYearlyPrices(res), nil
}

func (s *orderStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
