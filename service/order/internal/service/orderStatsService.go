package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-order/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/internal/repository"
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

type orderStatsService struct {
	errorhandler         errorhandler.OrderStatsError
	mencache             mencache.OrderStatsCache
	trace                trace.Tracer
	orderStatsRepository repository.OrderStatsRepository
	logger               logger.LoggerInterface
	mapping              response_service.OrderResponseMapper
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewOrderStatsService(
	errorhandler errorhandler.OrderStatsError,
	mencache mencache.OrderStatsCache,
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
		errorhandler:         errorhandler,
		mencache:             mencache,
		trace:                otel.Tracer("order-stats-service"),
		orderStatsRepository: orderStatsRepository,
		logger:               logger,
		mapping:              mapping,
		requestCounter:       requestCounter,
		requestDuration:      requestDuration,
	}
}

func (s *orderStatsService) FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTotalRevenue"

	year := req.Year
	month := req.Month

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalRevenueCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total revenue from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.orderStatsRepository.GetMonthlyTotalRevenue(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTotalRevenueError(err, method, "FAILED_FIND_MONTHLY_TOTAL_REVENUE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderMonthlyTotalRevenues(res)

	s.mencache.SetMonthlyTotalRevenueCache(ctx, req, so)

	logSuccess("Successfully fetched monthly total revenue", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *orderStatsService) FindYearlyTotalRevenue(ctx context.Context, year int) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	const method = "FindYearlyTotalRevenue"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalRevenueCache(ctx, year); found {
		logSuccess("Successfully fetched yearly total revenue from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.orderStatsRepository.GetYearlyTotalRevenue(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearlyTotalRevenueError(err, method, "FAILED_FIND_YEARLY_TOTAL_REVENUE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderYearlyTotalRevenues(res)

	s.mencache.SetYearlyTotalRevenueCache(ctx, year, so)

	logSuccess("Successfully fetched yearly total revenue", zap.Int("year", year))

	return so, nil
}

func (s *orderStatsService) FindMonthlyOrder(ctx context.Context, year int) ([]*response.OrderMonthlyResponse, *response.ErrorResponse) {
	const method = "FindMonthlyOrder"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyOrderCache(ctx, year); found {
		logSuccess("Successfully fetched monthly orders from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.orderStatsRepository.GetMonthlyOrder(ctx, year)

	if err != nil {
		return s.errorhandler.HandleMonthOrderStatsError(err, method, "FAILED_FIND_MONTHLY_ORDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderMonthlyPrices(res)

	s.mencache.SetMonthlyOrderCache(ctx, year, so)

	logSuccess("Successfully fetched monthly orders", zap.Int("year", year))

	return so, nil
}

func (s *orderStatsService) FindYearlyOrder(ctx context.Context, year int) ([]*response.OrderYearlyResponse, *response.ErrorResponse) {
	const method = "FindYearlyOrder"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyOrderCache(ctx, year); found {
		logSuccess("Successfully fetched yearly orders from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.orderStatsRepository.GetYearlyOrder(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearOrderStatsError(err, method, "FAILED_FIND_YEARLY_ORDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderYearlyPrices(res)

	s.mencache.SetYearlyOrderCache(ctx, year, so)

	logSuccess("Successfully fetched yearly orders", zap.Int("year", year))

	return so, nil
}

func (s *orderStatsService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

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

	return ctx, span, end, status, logSuccess
}

func (s *orderStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
