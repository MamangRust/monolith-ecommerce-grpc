package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
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

type merchantAwardHandleApi struct {
	client          pb.MerchantAwardServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantAwardResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerMerchantAward(
	router *echo.Echo,
	client pb.MerchantAwardServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantAwardResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
) *merchantAwardHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_award_handler_requests_total",
			Help: "Total number of banner requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_award_handler_request_duration_seconds",
			Help:    "Duration of banner requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantAwardHandler := &merchantAwardHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
		trace:           otel.Tracer("merchant-award-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/merchant-certification")

	routercategory.GET("", merchantAwardHandler.FindAllMerchantAward)
	routercategory.GET("/:id", merchantAwardHandler.FindById)
	routercategory.GET("/active", merchantAwardHandler.FindByActive)
	routercategory.GET("/trashed", merchantAwardHandler.FindByTrashed)

	routercategory.POST("/create", merchantAwardHandler.Create)
	routercategory.POST("/update/:id", merchantAwardHandler.Update)

	routercategory.POST("/trashed/:id", merchantAwardHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", merchantAwardHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", merchantAwardHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", merchantAwardHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", merchantAwardHandler.DeleteAllMerchantPermanent)

	return merchantAwardHandler
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
// @Success 200 {object} response.ApiResponsePaginationMerchantAward "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-certification [get]
func (h *merchantAwardHandleApi) FindAllMerchantAward(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllMerchantAward"
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

		logError("Failed to find all merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedFindAllMerchantAward(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantAward(res)

	logSuccess("Successfully find all merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantCertification
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAward "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-certification/{id} [get]
func (h *merchantAwardHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	req := &pb.FindByIdMerchantAwardRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to find merchant by ID", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedFindMerchantAwardById(c)
	}

	so := h.mapping.ToApiResponseMerchantAward(res)

	logSuccess("Successfully find merchant by ID", zap.Bool("success", true))

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
// @Success 200 {object} response.ApiResponsePaginationMerchantAwardDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-certification/active [get]
func (h *merchantAwardHandleApi) FindByActive(c echo.Context) error {
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

		logError("Failed to find active merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedFindActiveMerchantAward(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantAwardDeleteAt(res)

	logSuccess("Successfully find active merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantAwardDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-certification/trashed [get]
func (h *merchantAwardHandleApi) FindByTrashed(c echo.Context) error {
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

		logError("Failed to find trashed merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedFindTrashedMerchantAward(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantAwardDeleteAt(res)

	logSuccess("Successfully find trashed merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new merchant certification or award.
// @Summary Create a new merchant certification or award
// @Tags MerchantCertificationCertification
// @Description Create a new merchant certification or award with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantCertificationOrAwardRequest true "Create merchant certification or award request"
// @Success 200 {object} response.ApiResponseMerchantAward "Successfully created merchant certification or award"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant certification or award"
// @Router /api/merchant-certification/create [post]
func (h *merchantAwardHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	var body requests.CreateMerchantCertificationOrAwardRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Invalid request body", err, zap.Error(err))

		return merchantaward_errors.ErrApiBindCreateMerchantAward(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Invalid request body", err, zap.Error(err))

		return merchantaward_errors.ErrApiValidateCreateMerchantAward(c)
	}

	req := &pb.CreateMerchantAwardRequest{
		MerchantId:     int32(body.MerchantID),
		Title:          body.Title,
		Description:    body.Description,
		IssuedBy:       body.IssuedBy,
		IssueDate:      body.IssueDate,
		ExpiryDate:     body.ExpiryDate,
		CertificateUrl: body.CertificateUrl,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to create merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedCreateMerchantAward(c)
	}

	so := h.mapping.ToApiResponseMerchantAward(res)

	logSuccess("Successfully create merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing merchant certification or award.
// @Summary Update an existing merchant certification or award
// @Tags MerchantCertification
// @Description Update an existing merchant certification or award with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Merchant Certification ID"
// @Param request body requests.UpdateMerchantCertificationOrAwardRequest true "Update merchant certification or award request"
// @Success 200 {object} response.ApiResponseMerchantAward "Successfully updated merchant certification or award"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant certification or award"
// @Router /api/merchant-certification/update/{id} [post]
func (h *merchantAwardHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"

		logError("Invalid merchant award ID", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	var body requests.UpdateMerchantCertificationOrAwardRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Invalid request format", err, zap.Error(err))

		return merchantaward_errors.ErrApiBindUpdateMerchantAward(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Validation failed", err, zap.Error(err))

		return merchantaward_errors.ErrApiValidateUpdateMerchantAward(c)
	}

	req := &pb.UpdateMerchantAwardRequest{
		MerchantCertificationId: int32(idInt),
		Title:                   body.Title,
		Description:             body.Description,
		IssuedBy:                body.IssuedBy,
		IssueDate:               body.IssueDate,
		ExpiryDate:              body.ExpiryDate,
		CertificateUrl:          body.CertificateUrl,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		status = "error"

		logError("Merchant award update failed", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedUpdateMerchantAward(c)
	}

	so := h.mapping.ToApiResponseMerchantAward(res)

	logSuccess("Successfully updated merchant award", zap.Bool("success", true))

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
// @Success 200 {object} response.ApiResponseMerchantAwardDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant-certification/trashed/{id} [get]
func (h *merchantAwardHandleApi) TrashedMerchant(c echo.Context) error {
	const method = "TrashedMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("failed to archive merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	req := &pb.FindByIdMerchantAwardRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantAward(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to archive merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedTrashedMerchantAward(c)
	}

	so := h.mapping.ToApiResponseMerchantAwardDeleteAt(res)

	logSuccess("successfully retrieved trashed merchant", zap.Bool("success", true))

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
// @Success 200 {object} response.ApiResponseMerchantAwardDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-certification/restore/{id} [post]
func (h *merchantAwardHandleApi) RestoreMerchant(c echo.Context) error {
	const method = "RestoreMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID format", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	req := &pb.FindByIdMerchantAwardRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantAward(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedRestoreMerchantAward(c)
	}

	so := h.mapping.ToApiResponseMerchantAwardDeleteAt(res)

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
// @Router /api/merchant-certification/delete/{id} [delete]
func (h *merchantAwardHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	const method = "RestoreAllMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID format", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	req := &pb.FindByIdMerchantAwardRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantAwardPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to permanently delete merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedDeleteMerchantAwardPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantDelete(res)

	logSuccess("Successfully deleted merchant permanently", zap.Bool("success", true))

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
// @Router /api/merchant-certification/restore/all [post]
func (h *merchantAwardHandleApi) RestoreAllMerchant(c echo.Context) error {
	const method = "RestoreAllMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllMerchantAward(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all merchant", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedRestoreAllMerchantAward(c)
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
// @Router /api/merchant-certification/delete/all [post]
func (h *merchantAwardHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	const method = "DeleteAllMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllMerchantAwardPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Bulk merchant deletion failed", err, zap.Error(err))

		return merchantaward_errors.ErrApiFailedDeleteAllMerchantAwardPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	logSuccess("Successfully deleted all merchant permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantAwardHandleApi) startTracingAndLogging(
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

func (s *merchantAwardHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
