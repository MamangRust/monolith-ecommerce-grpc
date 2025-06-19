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
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderQueryService struct {
	ctx                  context.Context
	errorhandler         errorhandler.OrderQueryError
	mencache             mencache.OrderQueryCache
	trace                trace.Tracer
	orderQueryRepository repository.OrderQueryRepository
	logger               logger.LoggerInterface
	mapping              response_service.OrderResponseMapper
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewOrderQueryService(
	ctx context.Context,
	errorhandler errorhandler.OrderQueryError,
	mencache mencache.OrderQueryCache,
	orderQueryRepository repository.OrderQueryRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderResponseMapper,
) *orderQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_query_service_request_count",
			Help: "Total number of requests to the OrderQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_query_service_request_duration",
			Help:    "Histogram of request durations for the OrderQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &orderQueryService{
		ctx:                  ctx,
		errorhandler:         errorhandler,
		mencache:             mencache,
		trace:                otel.Tracer("order-query-service"),
		orderQueryRepository: orderQueryRepository,
		logger:               logger,
		mapping:              mapping,
		requestCounter:       requestCounter,
		requestDuration:      requestDuration,
	}
}

func (s *orderQueryService) FindAll(req *requests.FindAllOrder) ([]*response.OrderResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetOrderAllCache(req); found {
		logSuccess("Successfully fetched order from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindAllOrders(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_ORDERS", span, &status, zap.Error(err))
	}

	orderResponse := s.mapping.ToOrdersResponse(orders)

	s.mencache.SetOrderAllCache(req, orderResponse, totalRecords)

	logSuccess("Successfully fetched order", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return orderResponse, totalRecords, nil
}

func (s *orderQueryService) FindById(order_id int) (*response.OrderResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("order.id", order_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedOrderCache(order_id); found {
		logSuccess("Successfully fetched order from cache", zap.Int("order.id", order_id))

		return data, nil
	}

	order, err := s.orderQueryRepository.FindById(order_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_ORDER_BY_ID", span, &status, order_errors.ErrFailedFindOrderById, zap.Int("order_id", order_id))
	}

	so := s.mapping.ToOrderResponse(order)

	s.mencache.SetCachedOrderCache(so)

	logSuccess("Successfully fetched order", zap.Int("order.id", order_id))

	return so, nil
}

func (s *orderQueryService) FindByActive(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetOrderActiveCache(req); found {
		logSuccess("Successfully fetched order from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ALL_ORDERS_ACTIVE", span, &status, order_errors.ErrFailedFindOrdersByActive, zap.Error(err))
	}

	orderResponse := s.mapping.ToOrdersResponseDeleteAt(orders)

	s.mencache.SetOrderActiveCache(req, orderResponse, totalRecords)

	logSuccess("Successfully fetched order", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return orderResponse, totalRecords, nil
}

func (s *orderQueryService) FindByTrashed(req *requests.FindAllOrder) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetOrderTrashedCache(req); found {
		logSuccess("Successfully fetched order from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ALL_ORDERS_TRASHED", span, &status, order_errors.ErrFailedFindOrdersByTrashed, zap.Error(err))
	}

	orderResponse := s.mapping.ToOrdersResponseDeleteAt(orders)

	s.mencache.SetOrderTrashedCache(req, orderResponse, totalRecords)

	logSuccess("Successfully fetched order", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return orderResponse, totalRecords, nil
}

func (s *orderQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *orderQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *orderQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
