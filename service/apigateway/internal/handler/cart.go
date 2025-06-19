package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
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

type cartHandleApi struct {
	client          pb.CartServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.CartResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerCart(
	router *echo.Echo,
	client pb.CartServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CartResponseMapper,
) *cartHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_handler_requests_total",
			Help: "Total number of cart requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cart_handler_request_duration_seconds",
			Help:    "Duration of cart requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	cartHandler := &cartHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("cart-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCart := router.Group("/api/cart")
	routerCart.GET("", cartHandler.FindAll)
	routerCart.POST("/create", cartHandler.Create)
	routerCart.DELETE("/:id", cartHandler.Delete)
	routerCart.POST("/delete-all", cartHandler.DeleteAll)

	return cartHandler
}

// @Security Bearer
// @Summary Find all carts
// @Tags Cart
// @Description Retrieve a list of all carts
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponseCartPagination "List of carts"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve cart data"
// @Router /api/cart [get]
func (h *cartHandleApi) FindAll(c echo.Context) error {
	const method = "FindAllCart"

	ctx := c.Request().Context()

	userID, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil || userID <= 0 {
		h.logger.Debug("Invalid user ID format", zap.Error(err), zap.String("user_id", c.QueryParam("user_id")))
		return cart_errors.ErrApiFailedInvalidUserId(c)
	}

	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("user_id", userID),
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"
	defer func() { end(status) }()

	req := &pb.FindAllCartRequest{
		UserId:   int32(userID),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		status = "error"
		logError("Failed to fetch cart details", err, zap.Int("user_id", userID))
		return cart_errors.ErrApiFailedFindAllCarts(c)
	}

	resp := h.mapping.ToApiResponseCartPagination(res)
	logSuccess("Successfully fetched cart details", zap.Int("user_id", userID))

	return c.JSON(http.StatusOK, resp)
}

// @Security Bearer
// @Summary Create a new cart
// @Tags Cart
// @Description Create a new cart item
// @Accept json
// @Produce json
// @Param body body requests.CreateCartRequest true "Cart creation data"
// @Success 200 {object} response.ApiResponseCart "Created cart details"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 502 {object} response.ErrorResponse "Failed to create cart"
// @Router /api/cart [post]
func (h *cartHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	var body requests.CreateCartRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind create cart request", err, zap.Error(err))

		return cart_errors.ErrApiBindCreateCart(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate create cart request", err, zap.Error(err))

		return cart_errors.ErrApiValidateCreateCart(c)
	}

	req := &pb.CreateCartRequest{
		Quantity:  int32(body.Quantity),
		ProductId: int32(body.ProductID),
		UserId:    int32(body.UserID),
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to create cart", err, zap.Error(err))

		return cart_errors.ErrApiFailedCreateCart(c)
	}

	so := h.mapping.ToApiResponseCart(res)

	logSuccess("Successfully created cart", zap.Bool("success", true))

	return c.JSON(http.StatusCreated, so)
}

// @Security Bearer
// @Summary Delete a cart
// @Tags Cart
// @Description Delete a cart by ID (data dikirim via body)
// @Accept json
// @Produce json
// @Param request body requests.DeleteCartRequest true "Delete Cart Request"
// @Success 200 {object} response.ApiResponseCartDelete "Successfully deleted cart"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid input"
// @Failure 500 {object} response.ErrorResponse "Failed to delete cart"
// @Router /api/cart/delete [delete]
func (h *cartHandleApi) Delete(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	var req requests.DeleteCartRequest

	if err := c.Bind(&req); err != nil {
		status = "error"

		logError("Failed to bind delete cart request", err, zap.Error(err))

		return cart_errors.ErrApiBindDeleteCart(c)
	}

	if err := req.Validate(); err != nil {
		status = "error"

		logError("Failed to validate delete cart request", err, zap.Error(err))

		return cart_errors.ErrApiValidateDeleteCart(c)
	}

	reqPb := &pb.DeleteCartRequest{
		CartId: int32(req.CartID),
		UserId: int32(req.UserID),
	}

	res, err := h.client.Delete(ctx, reqPb)

	if err != nil {
		status = "error"

		logError("Failed to delete cart", err, zap.Error(err))

		return cart_errors.ErrApiFailedDeleteCart(c)
	}

	so := h.mapping.ToApiResponseCartDelete(res)

	logSuccess("Successfully deleted cart", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Delete multiple carts
// @Tags Cart
// @Description Delete multiple carts by IDs
// @Accept json
// @Produce json
// @Param request body requests.DeleteCartRequest true "Cart IDs"
// @Success 200 {object} response.ApiResponseCartAll "Successfully deleted carts"
// @Failure 500 {object} response.ErrorResponse "Failed to delete carts"
// @Router /api/cart/delete-all [post]
func (h *cartHandleApi) DeleteAll(c echo.Context) error {
	const method = "DeleteAll"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	var req requests.DeleteAllCartRequest
	if err := c.Bind(&req); err != nil {
		status = "error"

		logError("Failed to bind delete all cart request", err, zap.Error(err))

		return cart_errors.ErrApiBindDeleteAllCart(c)
	}

	if err := req.Validate(); err != nil {
		status = "error"

		logError("Failed to validate delete all cart request", err, zap.Error(err))

		return cart_errors.ErrApiValidateDeleteAllCart(c)
	}

	cartIdsPb := make([]int32, len(req.CartIds))
	for i, id := range req.CartIds {
		cartIdsPb[i] = int32(id)
	}

	reqPb := &pb.DeleteAllCartRequest{
		UserId:  int32(req.UserID),
		CartIds: cartIdsPb,
	}

	res, err := h.client.DeleteAll(ctx, reqPb)

	if err != nil {
		status = "error"

		logError("Failed to delete all cart", err, zap.Error(err))

		return cart_errors.ErrApiFailedDeleteAllCarts(c)
	}

	so := h.mapping.ToApiResponseCartAll(res)

	logSuccess("Successfully deleted all cart", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *cartHandleApi) startTracingAndLogging(
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

func (s *cartHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
