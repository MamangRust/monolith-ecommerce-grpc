package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/repository"
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

type cartQueryService struct {
	errorhandler        errorhandler.CartQueryError
	mencache            mencache.CartQueryCache
	trace               trace.Tracer
	cardQueryRepository repository.CartQueryRepository
	mapping             response_service.CartResponseMapper
	logger              logger.LoggerInterface
	requestCounter      *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewCartQueryService(
	errorhandler errorhandler.CartQueryError,
	mencache mencache.CartQueryCache,
	cardQueryRepository repository.CartQueryRepository, logger logger.LoggerInterface, mapping response_service.CartResponseMapper) *cartQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_query_service_request_count",
			Help: "Total number of requests to the CartQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cart_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the CartQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cartQueryService{
		errorhandler:        errorhandler,
		mencache:            mencache,
		trace:               otel.Tracer("cart-query-service"),
		cardQueryRepository: cardQueryRepository,
		mapping:             mapping,
		logger:              logger,
		requestCounter:      requestCounter,
		requestDuration:     requestDuration,
	}
}

func (s *cartQueryService) FindAll(ctx context.Context, req *requests.FindAllCarts) ([]*response.CartResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCartsCache(ctx, req); found {
		logSuccess("Successfully fetched all Carts from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	cart, totalRecords, err := s.cardQueryRepository.FindCarts(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_CARTS", span, &status, zap.Error(err))
	}

	cartRes := s.mapping.ToCartsResponse(cart)

	s.mencache.SetCartsCache(ctx, req, cartRes, totalRecords)

	logSuccess("Successfully fetched all Carts", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return cartRes, totalRecords, nil
}

func (s *cartQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

	s.logger.Info("Start: " + method)

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

func (s *cartQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *cartQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
