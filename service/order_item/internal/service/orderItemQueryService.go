package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderItemQueryService struct {
	errorhandler    errorhandler.OrderItemQueryError
	mencache        mencache.OrderItemQueryCache
	trace           trace.Tracer
	repo            repository.OrderItemQueryRepository
	mapping         response_service.OrderItemResponseMapper
	logger          logger.LoggerInterface
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewOrderItemQueryService(
	errorhandler errorhandler.OrderItemQueryError,
	mencache mencache.OrderItemQueryCache,
	repo repository.OrderItemQueryRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderItemResponseMapper,
) *orderItemQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_item_query_service_request_count",
			Help: "Total number of requests to the OrderItemQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_item_query_service_request_duration",
			Help:    "Histogram of request durations for the OrderItemQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &orderItemQueryService{
		errorhandler:    errorhandler,
		mencache:        mencache,
		trace:           otel.Tracer("order-item-query-service"),
		repo:            repo,
		mapping:         mapping,
		logger:          logger,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *orderItemQueryService) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, *response.ErrorResponse) {
	const method = "FindAllOrderItems"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedOrderItemsAll(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	orderItems, totalRecords, err := s.repo.FindAllOrderItems(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FindAllOrderItems", span, &status, zap.Error(err))
	}
	so := s.mapping.ToOrderItemsResponse(orderItems)

	s.mencache.SetCachedOrderItemsAll(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched order-items", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("totalRecords", *totalRecords))

	return so, totalRecords, nil
}

func (s *orderItemQueryService) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedOrderItemActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	orderItems, totalRecords, err := s.repo.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_ACTIVE_ORDER_ITEMS", span, &status, orderitem_errors.ErrFailedFindOrderItemsByActive, zap.Error(err))
	}

	so := s.mapping.ToOrderItemsResponseDeleteAt(orderItems)

	s.mencache.SetCachedOrderItemActive(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched order-items", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("totalRecords", *totalRecords))

	return so, totalRecords, nil
}

func (s *orderItemQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedOrderItemTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	orderItems, totalRecords, err := s.repo.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FindByTrashed", span, &status, orderitem_errors.ErrFailedFindOrderItemsByTrashed, zap.Error(err))
	}

	so := s.mapping.ToOrderItemsResponseDeleteAt(orderItems)

	s.mencache.SetCachedOrderItemTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched order-items", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("totalRecords", *totalRecords))

	return so, totalRecords, nil
}

func (s *orderItemQueryService) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse) {
	const method = "FindOrderItemByOrder"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedOrderItems(ctx, orderID); found {
		logSuccess("Successfully fetched order items from cache", zap.Int("order.id", orderID))

		return data, nil
	}

	orderItems, err := s.repo.FindOrderItemByOrder(ctx, orderID)

	if err != nil {
		return s.errorhandler.HandleRepositoryListError(err, method, "FindOrderItemByOrder", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	so := s.mapping.ToOrderItemsResponse(orderItems)

	s.mencache.SetCachedOrderItems(ctx, so)

	logSuccess("Successfully fetched order items", zap.Int("order.id", orderID))

	return so, nil
}

func (s *orderItemQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *orderItemQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *orderItemQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
