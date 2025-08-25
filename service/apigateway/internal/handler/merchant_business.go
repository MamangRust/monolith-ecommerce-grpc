package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
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

type merchantBusinessHandleApi struct {
	client          pb.MerchantBusinessServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantBusinessResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerMerchantBusiness(
	router *echo.Echo,
	client pb.MerchantBusinessServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantBusinessResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
) *merchantBusinessHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_business_handler_requests_total",
			Help: "Total number of merchant_business requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_business_handler_request_duration_seconds",
			Help:    "Duration of merchant_business requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantBusinessHandler := &merchantBusinessHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
		trace:           otel.Tracer("merchant-business-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/merchant-business")

	routercategory.GET("", merchantBusinessHandler.FindAllMerchantBusiness)
	routercategory.GET("/:id", merchantBusinessHandler.FindById)
	routercategory.GET("/active", merchantBusinessHandler.FindByActive)
	routercategory.GET("/trashed", merchantBusinessHandler.FindByTrashed)

	routercategory.POST("/create", merchantBusinessHandler.Create)
	routercategory.POST("/update/:id", merchantBusinessHandler.Update)

	routercategory.POST("/trashed/:id", merchantBusinessHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", merchantBusinessHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", merchantBusinessHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", merchantBusinessHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", merchantBusinessHandler.DeleteAllMerchantPermanent)

	return merchantBusinessHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags MerchantCertification
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusiness "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business [get]
func (h *merchantBusinessHandleApi) FindAllMerchantBusiness(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllMerchantBusiness"
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
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedFindAllMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantBusiness(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantCertification
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusiness "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/{id} [get]
func (h *merchantBusinessHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedFindByIdMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponseMerchantBusiness(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags MerchantCertification
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusinessDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/active [get]
func (h *merchantBusinessHandleApi) FindByActive(c echo.Context) error {
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
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedFindByActiveMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantBusinessDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantBusinessDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/trashed [get]
func (h *merchantBusinessHandleApi) FindByTrashed(c echo.Context) error {
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
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedFindByTrashedMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantBusinessDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new merchant business information
// @Tags MerchantBusiness
// @Description Create merchant business info (e.g., type, tax ID, website, etc.)
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantBusinessInformationRequest true "Create merchant business request"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Successfully created merchant business info"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant business info"
// @Router /api/merchant-business/create [post]
func (h *merchantBusinessHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	var body requests.CreateMerchantBusinessInformationRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Invalid request body", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiBindCreateMerchantBusiness(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Invalid request body", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiValidateCreateMerchantBusiness(c)
	}

	req := &pb.CreateMerchantBusinessRequest{
		MerchantId:        int32(body.MerchantID),
		BusinessType:      body.BusinessType,
		TaxId:             body.TaxID,
		EstablishedYear:   int32(body.EstablishedYear),
		NumberOfEmployees: int32(body.NumberOfEmployees),
		WebsiteUrl:        body.WebsiteUrl,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to create merchant business", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedCreateMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponseMerchantBusiness(res)

	logSuccess("Successfully created merchant business", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update existing merchant business information
// @Tags MerchantBusiness
// @Description Update merchant business info by ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant Business Info ID"
// @Param request body requests.UpdateMerchantBusinessInformationRequest true "Update merchant business request"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Successfully updated merchant business info"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant business info"
// @Router /api/merchant-business/update/{id} [post]
func (h *merchantBusinessHandleApi) Update(c echo.Context) error {
	const method = "DeleteAllBannerPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		status = "error"

		logError("Invalid merchant business ID parameter", err, zap.String("param_id", c.Param("id")))

		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	var body requests.UpdateMerchantBusinessInformationRequest
	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Invalid request body", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiBindUpdateMerchantBusiness(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Invalid request body", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiValidateUpdateMerchantBusiness(c)
	}

	req := &pb.UpdateMerchantBusinessRequest{
		MerchantBusinessInfoId: int32(idInt),
		BusinessType:           body.BusinessType,
		TaxId:                  body.TaxID,
		EstablishedYear:        int32(body.EstablishedYear),
		NumberOfEmployees:      int32(body.NumberOfEmployees),
		WebsiteUrl:             body.WebsiteUrl,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to update merchant business", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedUpdateMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponseMerchantBusiness(res)

	logSuccess("Successfully updated merchant business", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusinessDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant-business/trashed/{id} [get]
func (h *merchantBusinessHandleApi) TrashedMerchant(c echo.Context) error {
	const method = "TrashedMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID parameter", err, zap.String("param_id", c.Param("id")))

		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantBusiness(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve trashed merchant", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedTrashMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponseMerchantBusinessDeleteAt(res)

	logSuccess("Successfully retrieved trashed merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantCertification
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusinessDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-business/restore/{id} [post]
func (h *merchantBusinessHandleApi) RestoreMerchant(c echo.Context) error {
	const method = "RestoreMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID parameter", err, zap.String("param_id", c.Param("id")))

		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantBusiness(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore merchant", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedRestoreMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponseMerchantBusinessDeleteAt(res)

	logSuccess("Successfully restored merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantCertification
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-business/delete/{id} [delete]
func (h *merchantBusinessHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	const method = "DeleteMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID parameter", err, zap.String("param_id", c.Param("id")))

		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantBusinessPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to delete merchant", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedDeleteMerchantBusinessPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantDelete(res)

	logSuccess("Successfully deleted merchant record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantCertification
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-business/restore/all [post]
func (h *merchantBusinessHandleApi) RestoreAllMerchant(c echo.Context) error {
	const method = "RestoreAllMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllMerchantBusiness(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all merchant", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedRestoreAllMerchantBusiness(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	logSuccess("Successfully restored all merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantCertification
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-business/delete/all [post]
func (h *merchantBusinessHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	const method = "DeleteAllMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllMerchantBusinessPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to delete all merchant permanently", err, zap.Error(err))

		return merchantbusiness_errors.ErrApiFailedDeleteAllMerchantBusinessPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	logSuccess("Successfully deleted all merchant record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantBusinessHandleApi) startTracingAndLogging(
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

func (s *merchantBusinessHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
