package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderItemHandleApi struct {
	client          pb.OrderItemServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.OrderItemResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerOrderItem(
	router *echo.Echo,
	client pb.OrderItemServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.OrderItemResponseMapper,
) *orderItemHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_item_handler_requests_total",
			Help: "Total number of order item requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_item_handler_request_duration_seconds",
			Help:    "Duration of order item requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	categoryHandler := &orderItemHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("order-item-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/order-item")

	routercategory.GET("", categoryHandler.FindAllOrderItems)
	routercategory.GET("/:order_id", categoryHandler.FindOrderItemByOrder)
	routercategory.GET("/active", categoryHandler.FindByActive)
	routercategory.GET("/trashed", categoryHandler.FindByTrashed)

	return categoryHandler
}

// @Security Bearer
// @Summary Find all order items
// @Tags OrderItem
// @Description Retrieve a list of all order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItem "List of order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item [get]
func (h *orderItemHandleApi) FindAllOrderItems(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllOrderItems"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve order item data", err, zap.Error(err))

		return orderitem_errors.ErrApiOrderItemFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItem(res)

	logSuccess("successfully retrieve order item data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active order items
// @Tags OrderItem
// @Description Retrieve a list of active order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItemDeleteAt "List of active order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/active [get]
func (h *orderItemHandleApi) FindByActive(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByActive"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve order item data", err, zap.Error(err))

		return orderitem_errors.ErrApiOrderItemFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItemDeleteAt(res)

	logSuccess("successfully retrieve order item data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve trashed order items
// @Tags OrderItem
// @Description Retrieve a list of trashed order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItemDeleteAt "List of trashed order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/trashed [get]
func (h *orderItemHandleApi) FindByTrashed(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByTrashed"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve order item data", err, zap.Error(err))

		return orderitem_errors.ErrApiOrderItemFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItemDeleteAt(res)

	logSuccess("successfully retrieve order item data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find order items by order ID
// @Tags OrderItem
// @Description Retrieve order items by order ID
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.ApiResponsesOrderItem "List of order items by order ID"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/order/{order_id} [get]
func (h *orderItemHandleApi) FindOrderItemByOrder(c echo.Context) error {
	const method = "FindOrderItemByOrder"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	orderID, err := strconv.Atoi(c.Param("order_id"))

	if err != nil {
		status = "error"

		logError("Invalid order ID", err, zap.Error(err))

		return orderitem_errors.ErrApiOrderItemInvalidId(c)
	}

	req := &pb.FindByIdOrderItemRequest{
		Id: int32(orderID),
	}

	res, err := h.client.FindOrderItemByOrder(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to fetch order-items", err, zap.Error(err))

		return orderitem_errors.ErrApiOrderItemFailedFindByOrderId(c)
	}

	so := h.mapping.ToApiResponsesOrderItem(res)

	logSuccess("Successfully fetched order-items", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

func (s *orderItemHandleApi) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (func(string), func(string, ...zap.Field), func(string, error, ...zap.Field)) {
	start := time.Now()
	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := otelcode.Ok
		if status != "success" {
			code = otelcode.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError := func(msg string, err error, fields ...zap.Field) {
		span.RecordError(err)
		span.SetStatus(otelcode.Error, msg)
		span.AddEvent(msg)
		allFields := append([]zap.Field{zap.Error(err)}, fields...)
		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, logError
}

func (s *orderItemHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
