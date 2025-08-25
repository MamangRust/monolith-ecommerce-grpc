package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
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

type bannerHandleApi struct {
	client          pb.BannerServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.BannerResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandleBanner(
	router *echo.Echo,
	client pb.BannerServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.BannerResponseMapper,
) *bannerHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "banner_handler_requests_total",
			Help: "Total number of banner requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "banner_handler_request_duration_seconds",
			Help:    "Duration of banner requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	bannerHandler := &bannerHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("banner-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/banner")

	routercategory.GET("", bannerHandler.FindAllBanner)
	routercategory.GET("/:id", bannerHandler.FindById)
	routercategory.GET("/active", bannerHandler.FindByActive)
	routercategory.GET("/trashed", bannerHandler.FindByTrashed)

	routercategory.POST("/create", bannerHandler.Create)
	routercategory.POST("/update/:id", bannerHandler.Update)

	routercategory.POST("/trashed/:id", bannerHandler.TrashedBanner)
	routercategory.POST("/restore/:id", bannerHandler.RestoreBanner)
	routercategory.DELETE("/permanent/:id", bannerHandler.DeleteBannerPermanent)

	routercategory.POST("/restore/all", bannerHandler.RestoreAllBanner)
	routercategory.POST("/permanent/all", bannerHandler.DeleteAllBannerPermanent)

	return bannerHandler
}

// @Security Bearer
// @Summary Find all banners
// @Tags Banner
// @Description Retrieve a list of all banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBanner "List of banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner [get]
func (h *bannerHandleApi) FindAllBanner(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllBanner"
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

	req := &pb.FindAllBannerRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		status = "error"
		logError("Failed to find all banners", err, zap.Error(err))
		return banner_errors.ErrApiFailedFindAllBanner(c)
	}

	response := h.mapping.ToApiResponsePaginationBanner(res)

	logSuccess("Successfully retrieved all banners", zap.Bool("success", true))

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Find banner by ID
// @Tags Banner
// @Description Retrieve a banner by ID
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBanner "Banner data"
// @Failure 400 {object} response.ErrorResponse "Invalid banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/{id} [get]
func (h *bannerHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.String("id", c.Param("id")),
	)
	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid banner ID format", err, zap.Error(err))

		return banner_errors.ErrApiBannerInvalidId(c)
	}

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to find banner by id", err, zap.Error(err))
		return banner_errors.ErrApiFailedFindById(c)
	}

	so := h.mapping.ToApiResponseBanner(res)

	logSuccess("Successfully find banner by id", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active banners
// @Tags Banner
// @Description Retrieve a list of active banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBannerDeleteAt "List of active banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/active [get]
func (h *bannerHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllBannerRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to find active banners", err, zap.Error(err))

		return banner_errors.ErrApiFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationBannerDeleteAt(res)

	logSuccess("Successfully find active banners", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve trashed banners
// @Tags Banner
// @Description Retrieve a list of trashed banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBannerDeleteAt "List of active banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/trashed [get]
func (h *bannerHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllBannerRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to find trashed banners", err, zap.Error(err))

		return banner_errors.ErrApiFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationBannerDeleteAt(res)

	logSuccess("Successfully find trashed banners", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new banner.
// @Summary Create a new banner
// @Tags Banner
// @Description Create a new banner with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateBannerRequest true "Create banner request"
// @Success 200 {object} response.ApiResponseBanner "Successfully created banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create banner"
// @Router /api/banner/create [post]
func (h *bannerHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	var body requests.CreateBannerRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind create banner request", err, zap.Error(err))

		return banner_errors.ErrApiBindCreateBanner(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate create banner request", err, zap.Error(err))

		return banner_errors.ErrApiValidateCreateBanner(c)
	}

	req := &pb.CreateBannerRequest{
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		IsActive:  body.IsActive,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to create banner", err, zap.Error(err))

		return banner_errors.ErrApiFailedCreateBanner(c)
	}

	so := h.mapping.ToApiResponseBanner(res)

	logSuccess("Successfully create banner", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing banner record.
// @Summary Update an existing banner
// @Tags Banner
// @Description Update an existing banner record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Param request body requests.UpdateBannerRequest true "Update banner request"
// @Success 200 {object} response.ApiResponseBanner "Successfully updated banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update banner"
// @Router /api/banner/update/{id} [post]
func (h *bannerHandleApi) Update(c echo.Context) error {
	const method = "UpdateBanner"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		status = "error"
		logError("Invalid banner ID parameter", err, zap.String("param_id", c.Param("id")))
		return banner_errors.ErrApiBannerInvalidId(c)
	}

	var body requests.UpdateBannerRequest

	if err := c.Bind(&body); err != nil {
		status = "error"
		logError("Failed to bind UpdateBannerRequest", err, zap.String("param_id", c.Param("id")))
		return banner_errors.ErrApiBindUpdateBanner(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"
		logError("Validation failed for UpdateBannerRequest", err, zap.Any("body", body))
		return banner_errors.ErrApiValidateCreateBanner(c)
	}

	req := &pb.UpdateBannerRequest{
		BannerId:  int32(id),
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		IsActive:  body.IsActive,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		status = "error"
		logError("Failed to update banner", err, zap.Int("banner_id", id))
		return banner_errors.ErrApiFailedUpdateBanner(c)
	}

	resp := h.mapping.ToApiResponseBanner(res)
	logSuccess("Successfully updated banner", zap.Int("banner_id", id))

	return c.JSON(http.StatusOK, resp)
}

// @Security Bearer
// TrashedBanner retrieves a trashed Banner record by its ID.
// @Summary Retrieve a trashed Banner
// @Tags Banner
// @Description Retrieve a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDeleteAt "Successfully retrieved trashed Banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed Banner"
// @Router /api/banner/trashed/{id} [get]
func (h *bannerHandleApi) TrashedBanner(c echo.Context) error {
	const method = "TrashedBanner"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid banner ID parameter", err, zap.String("param_id", c.Param("id")))
		return banner_errors.ErrApiBannerInvalidId(c)
	}

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedBanner(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve trashed banner", err, zap.Int("banner.id", id))
		return banner_errors.ErrApiFailedTrashedBanner(c)
	}

	so := h.mapping.ToApiResponseBannerDeleteAt(res)

	logSuccess("Successfully retrieved trashed banner", zap.Int("banner.id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreBanner restores a Banner record from the trash by its ID.
// @Summary Restore a trashed Banner
// @Tags Banner
// @Description Restore a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDeleteAt "Successfully restored Banner"
// @Failure 400 {object} response.ErrorResponse "Invalid Banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore Banner"
// @Router /api/banner/restore/{id} [post]
func (h *bannerHandleApi) RestoreBanner(c echo.Context) error {
	const method = "RestoreBanner"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid banner ID parameter", err, zap.String("param_id", c.Param("id")))

		return banner_errors.ErrApiBannerInvalidId(c)
	}

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreBanner(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore banner", err, zap.Int("banner.id", id))

		return banner_errors.ErrApiFailedRestoreBanner(c)
	}

	so := h.mapping.ToApiResponseBannerDeleteAt(res)

	logSuccess("Successfully restored banner", zap.Int("banner.id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteBannerPermanent permanently deletes a Banner record by its ID.
// @Summary Permanently delete a Banner
// @Tags Banner
// @Description Permanently delete a Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDelete "Successfully deleted Banner record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete Banner:"
// @Router /api/banner/delete/{id} [delete]
func (h *bannerHandleApi) DeleteBannerPermanent(c echo.Context) error {
	const method = "DeleteBannerPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid banner ID parameter", err, zap.String("param_id", c.Param("id")))

		return banner_errors.ErrApiBannerInvalidId(c)
	}

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteBannerPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to delete banner permanently", err, zap.Int("banner.id", id))

		return banner_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseBannerDelete(res)

	logSuccess("Successfully deleted banner permanently", zap.Int("banner.id", id))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllBanner restores a Banner record from the trash by its ID.
// @Summary Restore a trashed Banner
// @Tags Banner
// @Description Restore a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerAll "Successfully restored Banner all"
// @Failure 400 {object} response.ErrorResponse "Invalid Banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore Banner"
// @Router /api/banner/restore/all [post]
func (h *bannerHandleApi) RestoreAllBanner(c echo.Context) error {
	const method = "RestoreAllBanner"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllBanner(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all banner", err, zap.Error(err))

		return banner_errors.ErrApiFailedRestoreAllBanner(c)
	}

	so := h.mapping.ToApiResponseBannerAll(res)

	logSuccess("Successfully restored all banner", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllBannerPermanent permanently deletes a banner record by its ID.
// @Summary Permanently delete a banner
// @Tags Banner
// @Description Permanently delete a banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "banner ID"
// @Success 200 {object} response.ApiResponseBannerAll "Successfully deleted banner record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete banner:"
// @Router /api/banner/delete/all [post]
func (h *bannerHandleApi) DeleteAllBannerPermanent(c echo.Context) error {
	const method = "DeleteAllBannerPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllBannerPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to delete all banner permanently", err)

		return banner_errors.ErrApiFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseBannerAll(res)

	logSuccess("Successfully deleted all banner permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *bannerHandleApi) startTracingAndLogging(
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

func (s *bannerHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
