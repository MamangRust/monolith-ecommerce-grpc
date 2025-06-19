package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
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

type merchantHandleApi struct {
	client          pb.MerchantServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerMerchant(
	router *echo.Echo,
	client pb.MerchantServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantResponseMapper,
) *merchantHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_handler_requests_total",
			Help: "Total number of merchant requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_handler_request_duration_seconds",
			Help:    "Duration of merchant requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	merchantHandler := &merchantHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("merchant-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/merchant")

	routercategory.GET("", merchantHandler.FindAllMerchant)
	routercategory.GET("/:id", merchantHandler.FindById)
	routercategory.GET("/active", merchantHandler.FindByActive)
	routercategory.GET("/trashed", merchantHandler.FindByTrashed)

	routercategory.POST("/create", merchantHandler.Create)
	routercategory.POST("/update/:id", merchantHandler.Update)

	routercategory.POST("/trashed/:id", merchantHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", merchantHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", merchantHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", merchantHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", merchantHandler.DeleteAllMerchantPermanent)

	return merchantHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags Merchant
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchant "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant [get]
func (h *merchantHandleApi) FindAllMerchant(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllMerchant"
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "failed"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchant(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags Merchant
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant/{id} [get]
func (h *merchantHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "failed"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMerchantById(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags Merchant
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant/active [get]
func (h *merchantHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "failed"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindActiveMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags Merchant
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant/trashed [get]
func (h *merchantHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "failed"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindTrashedMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new merchant.
// @Summary Create a new merchant
// @Tags Merchant
// @Description Create a new merchant with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantRequest true "Create merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Successfully created merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant"
// @Router /api/merchant/create [post]
func (h *merchantHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	var body requests.CreateMerchantRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind request body", err, zap.Error(err))

		return merchant_errors.ErrApiBindCreateMerchant(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate request body", err, zap.Error(err))

		return merchant_errors.ErrApiValidateCreateMerchant(c)
	}

	req := &pb.CreateMerchantRequest{
		UserId:       int32(body.UserID),
		Name:         body.Name,
		Description:  body.Description,
		Address:      body.Address,
		ContactEmail: body.ContactEmail,
		ContactPhone: body.ContactPhone,
		Status:       body.Status,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to create merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedCreateMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Successfully create merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing merchant record.
// @Summary Update an existing merchant
// @Tags Merchant
// @Description Update an existing merchant record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Param request body requests.UpdateMerchantRequest true "Update merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Successfully updated merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant"
// @Router /api/merchant/update [post]
func (h *merchantHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"

		logError("Failed to parse merchant id", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidId(c)
	}

	var body requests.UpdateMerchantRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind request body", err, zap.Error(err))

		return merchant_errors.ErrApiBindUpdateMerchant(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate request body", err, zap.Error(err))

		return merchant_errors.ErrApiValidateUpdateMerchant(c)
	}

	req := &pb.UpdateMerchantRequest{
		MerchantId:   int32(idInt),
		UserId:       int32(body.UserID),
		Name:         body.Name,
		Description:  body.Description,
		Address:      body.Address,
		ContactEmail: body.ContactEmail,
		ContactPhone: body.ContactPhone,
		Status:       body.Status,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to update merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedUpdateMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Successfully update merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags Merchant
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant/trashed/{id} [get]
func (h *merchantHandleApi) TrashedMerchant(c echo.Context) error {
	const method = "TrashedMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to parse merchant id", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to trashed merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedTrashMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchantDeleteAt(res)

	logSuccess("Successfully trashed merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags Merchant
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant/restore/{id} [post]
func (h *merchantHandleApi) RestoreMerchant(c echo.Context) error {
	const method = "RestoreMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to parse merchant id", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedRestoreMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Successfully restore merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags Merchant
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant/delete/{id} [delete]
func (h *merchantHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	const method = "DeleteMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to parse merchant id", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to delete merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedDeleteMerchantPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantDelete(res)

	logSuccess("Successfully deleted merchant record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags Merchant
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant/restore/all [post]
func (h *merchantHandleApi) RestoreAllMerchant(c echo.Context) error {
	const method = "RestoreAllMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllMerchant(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedRestoreAllMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchantAll(res)

	logSuccess("Successfully restore all merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags Merchant
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant/delete/all [post]
func (h *merchantHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	const method = "DeleteAllMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllMerchantPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Bulk merchant deletion failed", err, zap.Error(err))

		return merchant_errors.ErrApiFailedDeleteAllMerchantPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantAll(res)

	logSuccess("Bulk merchant deletion success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantHandleApi) startTracingAndLogging(
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

func (s *merchantHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
