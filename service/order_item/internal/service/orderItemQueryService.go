package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
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
	ctx             context.Context
	trace           trace.Tracer
	repo            repository.OrderItemQueryRepository
	mapping         response_service.OrderItemResponseMapper
	logger          logger.LoggerInterface
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewOrderItemQueryService(
	ctx context.Context,
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
		ctx:             ctx,
		trace:           otel.Tracer("order-item-query-service"),
		repo:            repo,
		mapping:         mapping,
		logger:          logger,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *orderItemQueryService) FindAllOrderItems(req *requests.FindAllOrderItems) ([]*response.OrderItemResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllOrderItems", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllOrderItems")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all order items",
		zap.String("search", search),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orderItems, totalRecords, err := s.repo.FindAllOrderItems(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_ORDER_ITEM")

		s.logger.Error("Failed to retrieve order-items",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("page", page),
			attribute.Int("pageSize", pageSize),
			attribute.String("search", search),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order-items")

		status = "failed_find_all_order_item"

		return nil, nil, orderitem_errors.ErrFailedFindAllOrderItems
	}

	s.logger.Debug("Successfully fetched order-item",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToOrderItemsResponse(orderItems), totalRecords, nil
}

func (s *orderItemQueryService) FindByActive(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all order items active",
		zap.String("search", search),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orderItems, totalRecords, err := s.repo.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEM_BY_ACTIVE")

		s.logger.Error("Failed to retrieve order-items",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("page", page),
			attribute.Int("pageSize", pageSize),
			attribute.String("search", search),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order-items")

		status = "failed_find_order_item_by_active"

		return nil, nil, orderitem_errors.ErrFailedFindOrderItemsByActive
	}

	s.logger.Debug("Successfully fetched order-items",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToOrderItemsResponseDeleteAt(orderItems), totalRecords, nil
}

func (s *orderItemQueryService) FindByTrashed(req *requests.FindAllOrderItems) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all order items trashed",
		zap.String("search", search),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orderItems, totalRecords, err := s.repo.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEM_BY_TRASHED")

		s.logger.Error("Failed to retrieve order-items",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.Int("page", page),
			attribute.Int("pageSize", pageSize),
			attribute.String("search", search),
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order-items")

		status = "failed_find_order_item_by_trashed"

		return nil, nil, orderitem_errors.ErrFailedFindOrderItemsByTrashed
	}

	s.logger.Debug("Successfully fetched order-items",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToOrderItemsResponseDeleteAt(orderItems), totalRecords, nil
}

func (s *orderItemQueryService) FindOrderItemByOrder(orderID int) ([]*response.OrderItemResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching order items by order", zap.Int("order_id", orderID))

	orderItems, err := s.repo.FindOrderItemByOrder(orderID)

	if err != nil {
		s.logger.Error("Failed to retrieve order items",
			zap.Error(err),
			zap.Int("order_id", orderID))
		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}
	return s.mapping.ToOrderItemsResponse(orderItems), nil
}

func (s *orderItemQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
