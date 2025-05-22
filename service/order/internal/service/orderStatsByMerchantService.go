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

type orderStatsByMerchantService struct {
	ctx                            context.Context
	trace                          trace.Tracer
	orderStatsByMerchantRepository repository.OrderStatsByMerchantRepository
	mapping                        response_service.OrderResponseMapper
	logger                         logger.LoggerInterface
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewOrderStatsByMerchantService(
	ctx context.Context,
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
		ctx:                            ctx,
		trace:                          otel.Tracer("order-stats-by-merchant-service"),
		orderStatsByMerchantRepository: orderStatsByMerchantRepository,
		logger:                         logger,
		mapping:                        mapping,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *orderStatsByMerchantService) FindMonthlyTotalRevenueByMerchant(req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalRevenueByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalRevenueByMerchant")
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

	res, err := s.orderStatsByMerchantRepository.GetMonthlyTotalRevenueByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_REVENUE_BY_MERCHANT")

		s.logger.Error("failed to get monthly total revenue",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "failed to get monthly total revenue")

		status = "failed_find_monthly_total_revenue_by_merchant"

		return nil, order_errors.ErrFailedFindMonthlyTotalRevenueByMerchant
	}

	return s.mapping.ToOrderMonthlyTotalRevenues(res), nil
}

func (s *orderStatsByMerchantService) FindYearlyTotalRevenueByMerchant(req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalRevenueByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalRevenueByMerchant")
	defer span.End()

	year := req.Year
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant_id", merchantId),
	)

	s.logger.Debug("find yearly total revenue",
		zap.Int("year", year),
		zap.Int("merchant_id", merchantId))

	res, err := s.orderStatsByMerchantRepository.GetYearlyTotalRevenueByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_REVENUE_BY_MERCHANT")

		s.logger.Error("failed to get yearly total revenue",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)

		span.SetStatus(codes.Error, "failed to get yearly total revenue")

		status = "failed_find_yearly_total_revenue_by_merchant"

		return nil, order_errors.ErrFailedFindYearlyTotalRevenueByMerchant
	}

	return s.mapping.ToOrderYearlyTotalRevenues(res), nil
}

func (s *orderStatsByMerchantService) FindMonthlyOrderByMerchant(req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyOrderByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyOrderByMerchant")
	defer span.End()

	year := req.Year
	merchant_id := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("find monthly orders by merchant",
		zap.Int("year", year),
		zap.Int("merchant_id", merchant_id))

	res, err := s.orderStatsByMerchantRepository.GetMonthlyOrderByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_ORDER_BY_MERCHANT")

		s.logger.Error("failed to get monthly orders by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(
			codes.Error,
			"failed to get monthly orders by merchant",
		)

		status = "failed_find_monthly_order_by_merchant"

		return nil, order_errors.ErrFailedFindMonthlyOrderByMerchant
	}

	return s.mapping.ToOrderMonthlyPrices(res), nil
}

func (s *orderStatsByMerchantService) FindYearlyOrderByMerchant(req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyOrderByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyOrderByMerchant")
	defer span.End()

	year := req.Year
	merchant_id := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("find yearly orders by merchant",
		zap.Int("year", year),
		zap.Int("merchant_id", merchant_id))

	res, err := s.orderStatsByMerchantRepository.GetYearlyOrderByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_ORDER_BY_MERCHANT")

		s.logger.Error("failed to get yearly orders by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(
			codes.Error,
			"Failed to get yearly orders by merchant",
		)

		status = "failed_find_yearly_order_by_merchant"

		return nil, order_errors.ErrFailedFindYearlyOrderByMerchant
	}

	return s.mapping.ToOrderYearlyPrices(res), nil
}

func (s *orderStatsByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
