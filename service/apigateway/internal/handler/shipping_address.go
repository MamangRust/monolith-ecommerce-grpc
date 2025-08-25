package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingAddressHandleApi struct {
	client          pb.ShippingServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.ShippingAddressResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerShippingAddress(
	router *echo.Echo,
	client pb.ShippingServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ShippingAddressResponseMapper,
) *shippingAddressHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shipping_handler_requests_total",
			Help: "Total number of shipping requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shipping_handler_request_duration_seconds",
			Help:    "Duration of shipping requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	shippingHandler := &shippingAddressHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("shipping-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	myrouter := router.Group("/api/shipping-address")

	myrouter.GET("", shippingHandler.FindAllShipping)
	myrouter.GET("/:id", shippingHandler.FindById)
	myrouter.GET("/order/:id", shippingHandler.FindByOrder)
	myrouter.GET("/active", shippingHandler.FindByActive)
	myrouter.GET("/trashed", shippingHandler.FindByTrashed)

	myrouter.POST("/trashed/:id", shippingHandler.TrashedShippingAddress)
	myrouter.POST("/restore/:id", shippingHandler.RestoreShippingAddress)
	router.DELETE("/permanent/:id", shippingHandler.DeleteShippingAddressPermanent)

	myrouter.POST("/restore/all", shippingHandler.RestoreAllShippingAddress)
	myrouter.POST("/permanent/all", shippingHandler.DeleteAllShippingAddressPermanent)

	return shippingHandler

}

// @Security Bearer
// @Summary Find all shipping-address
// @Tags shipping address
// @Description Retrieve a list of all shipping-address
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddress "List of shipping-address"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping-address data"
// @Router /api/shipping-address [get]
func (h *shippingAddressHandleApi) FindAllShipping(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllShipping"
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

	req := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve shipping-address data", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedFindAllShippingAddresses(c)
	}

	so := h.mapping.ToApiResponsePaginationShippingAddress(res)

	logSuccess("Successfully retrieved shipping-address data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find shipping address by ID
// @Tags ShippingAddress
// @Description Retrieve a shipping address by ID
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddress "Shipping address data"
// @Failure 400 {object} response.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address/{id} [get]
func (h *shippingAddressHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid shipping address ID", err, zap.Error(err))

		return shippingaddress_errors.ErrApiInvalidIdShippingAddress(c)
	}

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve shipping address data", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedFindShippingAddressById(c)
	}

	so := h.mapping.ToApiResponseShippingAddress(res)

	logSuccess("Successfully retrieved shipping address data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find shipping address by order ID
// @Tags ShippingAddress
// @Description Retrieve a shipping address by order ID
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.ApiResponseShippingAddress "Shipping address data"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address/order/{order_id} [get]
func (h *shippingAddressHandleApi) FindByOrder(c echo.Context) error {
	const method = "FindByOrder"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	orderID, err := strconv.Atoi(c.Param("order_id"))

	if err != nil {
		status = "error"

		logError("Invalid order ID", err, zap.Error(err))

		return shippingaddress_errors.ErrApiInvalidOrderIdShippingAddress(c)
	}

	req := &pb.FindByIdShippingRequest{
		Id: int32(orderID),
	}

	res, err := h.client.FindByOrder(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve shipping address data", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedFindShippingAddressByOrder(c)
	}

	so := h.mapping.ToApiResponseShippingAddress(res)

	logSuccess("Successfully retrieved shipping address data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active shipping-address
// @Tags ShippingAddress
// @Description Retrieve a list of active shipping-address
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddressDeleteAt "List of active shipping-address"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping data"
// @Router /api/shipping-address/active [get]
func (h *shippingAddressHandleApi) FindByActive(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllActive"
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

	req := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to fetch shipping address details", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedFindActiveShippingAddresses(c)
	}

	so := h.mapping.ToApiResponsePaginationShippingAddressDeleteAt(res)

	logSuccess("Successfully fetch shipping address details", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed shipping-address records.
// @Summary Retrieve trashed shipping-address
// @Tags ShippingAddress
// @Description Retrieve a list of trashed shipping-address records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddressDeleteAt "List of trashed shipping-address data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping-address data"
// @Router /api/shipping-address/trashed [get]
func (h *shippingAddressHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to fetch shipping address details", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedFindTrashedShippingAddresses(c)
	}

	so := h.mapping.ToApiResponsePaginationShippingAddressDeleteAt(res)

	logSuccess("Successfully fetch shipping address details", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedShippingAddress retrieves a trashed shipping address record by its ID.
// @Summary Retrieve a trashed shipping address
// @Tags ShippingAddress
// @Description Retrieve a trashed shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDeleteAt "Successfully retrieved trashed shipping address"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed shipping address"
// @Router /api/shipping-address/trashed/{id} [get]
func (h *shippingAddressHandleApi) TrashedShippingAddress(c echo.Context) error {
	const method = "TrashedShippingAddress"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to fetch shipping address details", err, zap.Error(err))

		return shippingaddress_errors.ErrApiInvalidIdShippingAddress(c)
	}

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedShipping(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to fetch shipping address details", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedTrashShippingAddress(c)
	}

	so := h.mapping.ToApiResponseShippingAddressDeleteAt(res)

	logSuccess("Successfully fetch shipping address details", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreShippingAddress restores a shipping address record from the trash by its ID.
// @Summary Restore a trashed shipping address
// @Tags ShippingAddress
// @Description Restore a trashed shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDeleteAt "Successfully restored shipping address"
// @Failure 400 {object} response.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore shipping address"
// @Router /api/shipping-address/restore/{id} [post]
func (h *shippingAddressHandleApi) RestoreShippingAddress(c echo.Context) error {
	const method = "RestoreShippingAddress"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to restore shipping address", err, zap.Error(err))

		return shippingaddress_errors.ErrApiInvalidIdShippingAddress(c)
	}

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreShipping(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore shipping address", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedRestoreShippingAddress(c)
	}

	so := h.mapping.ToApiResponseShippingAddressDeleteAt(res)

	logSuccess("Successfully restore shipping address", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteShippingAddressPermanent permanently deletes a shipping address record by its ID.
// @Summary Permanently delete a shipping address
// @Tags ShippingAddress
// @Description Permanently delete a shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDelete "Successfully deleted shipping address record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete shipping address:"
// @Router /api/shipping-address/delete/{id} [delete]
func (h *shippingAddressHandleApi) DeleteShippingAddressPermanent(c echo.Context) error {
	const method = "DeleteShippingAddressPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to permanently delete shipping address", err, zap.Error(err))

		return shippingaddress_errors.ErrApiInvalidIdShippingAddress(c)
	}

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteShippingPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to permanently delete shipping address", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedDeleteShippingAddressPermanent(c)
	}

	so := h.mapping.ToApiResponseShippingAddressDelete(res)

	logSuccess("Successfully deleted shipping address record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllShippingAddress restores all trashed shipping address records.
// @Summary Restore all trashed shipping addresses
// @Tags ShippingAddress
// @Description Restore all trashed shipping address records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseShippingAddressAll "Successfully restored all shipping addresses"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all shipping addresses"
// @Router /api/shipping-address/restore/all [post]
func (h *shippingAddressHandleApi) RestoreAllShippingAddress(c echo.Context) error {
	const method = "RestoreAllShippingAddress"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllShipping(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all shipping addresses", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedRestoreAllShippingAddresses(c)
	}

	so := h.mapping.ToApiResponseShippingAddressAll(res)

	logSuccess("Successfully restored all shipping addresses", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllShippingAddressPermanent permanently deletes all trashed shipping address records.
// @Summary Permanently delete all trashed shipping addresses
// @Tags ShippingAddress
// @Description Permanently delete all trashed shipping address records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseShippingAddressAll "Successfully deleted all shipping addresses permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete all shipping addresses"
// @Router /api/shipping-address/delete/all [post]
func (h *shippingAddressHandleApi) DeleteAllShippingAddressPermanent(c echo.Context) error {
	const method = "DeleteAllShippingAddressPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllShippingPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to permanently delete all shipping addresses", err, zap.Error(err))

		return shippingaddress_errors.ErrApiFailedDeleteAllPermanentShippingAddresses(c)
	}

	so := h.mapping.ToApiResponseShippingAddressAll(res)

	logSuccess("Successfully deleted all shipping addresses permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *shippingAddressHandleApi) startTracingAndLogging(
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

func (s *shippingAddressHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
