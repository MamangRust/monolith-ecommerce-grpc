package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cartQueryService struct {
	ctx                 context.Context
	trace               trace.Tracer
	cardQueryRepository repository.CartQueryRepository
	mapping             response_service.CartResponseMapper
	logger              logger.LoggerInterface
	requestCounter      *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewCartQueryService(ctx context.Context, cardQueryRepository repository.CartQueryRepository, logger logger.LoggerInterface, mapping response_service.CartResponseMapper) *cartQueryService {
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
		ctx:                 ctx,
		trace:               otel.Tracer("cart-query-service"),
		cardQueryRepository: cardQueryRepository,
		mapping:             mapping,
		logger:              logger,
		requestCounter:      requestCounter,
		requestDuration:     requestDuration,
	}
}

func (s *cartQueryService) FindAll(req *requests.FindAllCarts) ([]*response.CartResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching cart",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	cart, totalRecords, err := s.cardQueryRepository.FindCarts(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_CARTS")

		s.logger.Error("Failed to fetch cart",
			zap.String("traceID", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to fetch cart")

		status = "failed_find_all_carts"

		return nil, nil, cart_errors.ErrFailedFindAllCarts
	}

	cartRes := s.mapping.ToCartsResponse(cart)

	s.logger.Debug("Successfully fetched cart",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return cartRes, totalRecords, nil
}

func (s *cartQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
