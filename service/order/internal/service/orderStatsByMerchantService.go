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

type orderStatsByMerchantService struct {
	mencache                       mencache.OrderStatsByMerchantCache
	errorhandler                   errorhandler.OrderStatsByMerchantError
	trace                          trace.Tracer
	orderStatsByMerchantRepository repository.OrderStatsByMerchantRepository
	mapping                        response_service.OrderResponseMapper
	logger                         logger.LoggerInterface
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewOrderStatsByMerchantService(
	mencache mencache.OrderStatsByMerchantCache,
	errorhandler errorhandler.OrderStatsByMerchantError,
	orderStatsByMerchantRepository repository.OrderStatsByMerchantRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderResponseMapper,
) *orderStatsByMerchantService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_stats_by_merchant_service_request_count",
			Help: "Total number of requests to the OrderStatsByMerchantService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_stats_by_merchant_service_request_duration",
			Help:    "Histogram of request durations for the OrderStatsByMerchantService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &orderStatsByMerchantService{
		mencache:                       mencache,
		errorhandler:                   errorhandler,
		trace:                          otel.Tracer("order-stats-by-merchant-service"),
		orderStatsByMerchantRepository: orderStatsByMerchantRepository,
		logger:                         logger,
		mapping:                        mapping,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *orderStatsByMerchantService) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTotalRevenueByMerchant"

	merchantId := req.MerchantID
	year := req.Year
	month := req.Month

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantId), attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalRevenueByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total revenue from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetMonthlyTotalRevenueByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthTotalRevenueByMerchantError(err, method, "FAILED_FIND_MONTHLY_TOTAL_REVENUE_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderMonthlyTotalRevenues(res)
	s.mencache.SetMonthlyTotalRevenueByMerchantCache(ctx, req, so)

	logSuccess("Successfully fetched monthly total revenue", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *orderStatsByMerchantService) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	const method = "FindYearlyTotalRevenueByMerchant"

	year := req.Year
	merchantId := req.MerchantID

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantId), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalRevenueByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total revenue from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetYearlyTotalRevenueByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleYearTotalRevenueByMerchantError(err, method, "FAILED_FIND_YEARLY_TOTAL_REVENUE_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderYearlyTotalRevenues(res)
	s.mencache.SetYearlyTotalRevenueByMerchantCache(ctx, req, so)

	logSuccess("Successfully fetched yearly total revenue", zap.Int("year", year))

	return so, nil
}

func (s *orderStatsByMerchantService) FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, *response.ErrorResponse) {
	const method = "FindMonthlyOrderByMerchant"

	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchant_id), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly orders from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetMonthlyOrderByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthOrderStatsByMerchantError(err, method, "FAILED_FIND_MONTHLY_ORDER_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderMonthlyPrices(res)
	s.mencache.SetMonthlyOrderByMerchantCache(ctx, req, so)

	logSuccess("Successfully fetched monthly orders", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

	return so, nil
}

func (s *orderStatsByMerchantService) FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, *response.ErrorResponse) {
	const method = "FindYearlyOrderByMerchant"

	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchant_id), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly orders from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

		return data, nil
	}

	res, err := s.orderStatsByMerchantRepository.GetYearlyOrderByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleYearOrderStatsByMerchantError(err, "FindYearlyOrderByMerchant", "FAILED_FIND_YEARLY_ORDER_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderYearlyPrices(res)
	s.mencache.SetYearlyOrderByMerchantCache(ctx, req, so)

	logSuccess("Successfully fetched yearly orders", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

	return so, nil
}

func (s *orderStatsByMerchantService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *orderStatsByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
