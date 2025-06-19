package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
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

type sliderHandleApi struct {
	client          pb.SliderServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.SliderResponseMapper
	upload_image    upload_image.ImageUploads
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerSlider(
	router *echo.Echo,
	client pb.SliderServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.SliderResponseMapper,
	upload_image upload_image.ImageUploads,
) *sliderHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "slider_handler_requests_total",
			Help: "Total number of slider requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "slider_handler_request_duration_seconds",
			Help:    "Duration of slider requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	sliderHandler := &sliderHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		upload_image:    upload_image,
		trace:           otel.Tracer("slider-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerSlider := router.Group("/api/slider")

	routerSlider.GET("", sliderHandler.FindAllSlider)
	routerSlider.GET("/active", sliderHandler.FindByActive)
	routerSlider.GET("/trashed", sliderHandler.FindByTrashed)

	routerSlider.POST("/create", sliderHandler.Create)
	routerSlider.POST("/update/:id", sliderHandler.Update)

	routerSlider.POST("/trashed/:id", sliderHandler.TrashedSlider)
	routerSlider.POST("/restore/:id", sliderHandler.RestoreSlider)
	routerSlider.DELETE("/permanent/:id", sliderHandler.DeleteSliderPermanent)

	routerSlider.POST("/restore/all", sliderHandler.RestoreAllSlider)
	routerSlider.POST("/permanent/all", sliderHandler.DeleteAllSliderPermanent)

	return sliderHandler

}

// @Security Bearer
// @Summary Find all slider
// @Tags Slider
// @Description Retrieve a list of all slider
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSlider "List of slider"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider [get]
func (h *sliderHandleApi) FindAllSlider(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllSlider"
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

	req := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve slider data", err, zap.Error(err))

		return slider_errors.ErrApiFailedFindAllSliders(c)
	}

	so := h.mapping.ToApiResponsePaginationSlider(res)

	logSuccess("Successfully retrieve slider data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active slider
// @Tags Slider
// @Description Retrieve a list of active slider
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSliderDeleteAt "List of active slider"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider/active [get]
func (h *sliderHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve slider data", err, zap.Error(err))

		return slider_errors.ErrApiFailedFindActiveSliders(c)
	}

	so := h.mapping.ToApiResponsePaginationSliderDeleteAt(res)

	logSuccess("Successfully retrieve slider data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed slider records.
// @Summary Retrieve trashed slider
// @Tags Slider
// @Description Retrieve a list of trashed slider records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSliderDeleteAt "List of trashed slider data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider/trashed [get]
func (h *sliderHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve trashed slider data", err, zap.Error(err))

		return slider_errors.ErrApiFailedFindTrashedSliders(c)
	}

	so := h.mapping.ToApiResponsePaginationSliderDeleteAt(res)

	logSuccess("Successfully retrieve trashed slider data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new slider with image upload.
// @Summary Create a new slider
// @Tags Slider
// @Description Create a new slider with the provided details and an image file
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Slider name"
// @Param image_slider formData file true "Slider image file"
// @Success 200 {object} response.ApiResponseSlider "Successfully created slider"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create slider"
// @Router /api/slider/create [post]
func (h *sliderHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	formData, err := h.parseSliderForm(c, true)
	if err != nil {
		status = "error"

		logError("Failed to create slider", err, zap.Error(err))

		return slider_errors.ErrApiFailedCreateSlider(c)
	}

	req := &pb.CreateSliderRequest{
		Name:  formData.Nama,
		Image: formData.FilePath,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to create slider", err, zap.Error(err))

		return slider_errors.ErrApiFailedCreateSlider(c)
	}

	so := h.mapping.ToApiResponseSlider(res)

	logSuccess("Successfully created slider", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing slider with image upload.
// @Summary Update an existing slider
// @Tags Slider
// @Description Update an existing slider record with the provided details and an optional image file
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Slider ID"
// @Param name formData string true "Slider name"
// @Param image_slider formData file false "New slider image file"
// @Success 200 {object} response.ApiResponseSlider "Successfully updated slider"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update slider"
// @Router /api/slider/update [post]
func (h *sliderHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	sliderID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to update slider", err, zap.Error(err))

		return slider_errors.ErrApiInvalidId(c)
	}

	formData, err := h.parseSliderForm(c, true)

	if err != nil {
		status = "error"

		logError("Failed to update slider", err, zap.Error(err))

		return slider_errors.ErrApiInvalidBody(c)
	}

	req := &pb.UpdateSliderRequest{
		Id:    int32(sliderID),
		Name:  formData.Nama,
		Image: formData.FilePath,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to update slider", err, zap.Error(err))

		return slider_errors.ErrApiFailedUpdateSlider(c)
	}

	so := h.mapping.ToApiResponseSlider(res)

	logSuccess("Successfully updated slider", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedSlider retrieves a trashed slider record by its ID.
// @Summary Retrieve a trashed slider
// @Tags Slider
// @Description Retrieve a trashed slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderDeleteAt "Successfully retrieved trashed slider"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed slider"
// @Router /api/slider/trashed/{id} [get]
func (h *sliderHandleApi) TrashedSlider(c echo.Context) error {
	const method = "TrashedSlider"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to archive slider", err, zap.Error(err))

		return slider_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedSlider(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to archive slider", err, zap.Error(err))

		return slider_errors.ErrApiFailedTrashSlider(c)
	}

	so := h.mapping.ToApiResponseSliderDeleteAt(res)

	logSuccess("Successfully trashed slider", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreSlider restores a slider record from the trash by its ID.
// @Summary Restore a trashed slider
// @Tags Slider
// @Description Restore a trashed slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderDeleteAt "Successfully restored slider"
// @Failure 400 {object} response.ErrorResponse "Invalid slider ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore slider"
// @Router /api/slider/restore/{id} [post]
func (h *sliderHandleApi) RestoreSlider(c echo.Context) error {
	const method = "RestoreSlider"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to restore slider", err, zap.Error(err))

		return slider_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreSlider(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore slider", err, zap.Error(err))

		return slider_errors.ErrApiFailedRestoreSlider(c)
	}

	so := h.mapping.ToApiResponseSliderDeleteAt(res)

	logSuccess("Successfully restored slider", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteSliderPermanent permanently deletes a slider record by its ID.
// @Summary Permanently delete a slider
// @Tags Slider
// @Description Permanently delete a slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderDelete "Successfully deleted slider record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete slider:"
// @Router /api/slider/delete/{id} [delete]
func (h *sliderHandleApi) DeleteSliderPermanent(c echo.Context) error {
	const method = "DeleteSliderPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to permanently delete slider", err, zap.Error(err))

		return slider_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdSliderRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteSliderPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to permanently delete slider", err, zap.Error(err))

		return slider_errors.ErrApiFailedDeleteSliderPermanent(c)
	}

	so := h.mapping.ToApiResponseSliderDelete(res)

	logSuccess("Successfully deleted slider record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllSlider restores a slider record from the trash by its ID.
// @Summary Restore a trashed slider
// @Tags Slider
// @Description Restore a trashed slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderAll "Successfully restored slider all"
// @Failure 400 {object} response.ErrorResponse "Invalid slider ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore slider"
// @Router /api/slider/restore/all [post]
func (h *sliderHandleApi) RestoreAllSlider(c echo.Context) error {
	const method = "RestoreAllSlider"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllSlider(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all slider", err, zap.Error(err))

		return slider_errors.ErrApiFailedRestoreAllSliders(c)
	}

	so := h.mapping.ToApiResponseSliderAll(res)

	logSuccess("Successfully restored all slider", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllSliderPermanent permanently deletes a slider record by its ID.
// @Summary Permanently delete a slider
// @Tags Slider
// @Description Permanently delete a slider record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "slider ID"
// @Success 200 {object} response.ApiResponseSliderAll "Successfully deleted slider record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete slider:"
// @Router /api/slider/delete/all [post]
func (h *sliderHandleApi) DeleteAllSliderPermanent(c echo.Context) error {
	const method = "DeleteAllSliderPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllSliderPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to permanently delete all slider", err, zap.Error(err))

		return slider_errors.ErrApiFailedDeleteAllPermanentSliders(c)
	}

	so := h.mapping.ToApiResponseSliderAll(res)

	logSuccess("Successfully deleted all slider record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (h *sliderHandleApi) parseSliderForm(c echo.Context, requireImage bool) (requests.SliderFormData, error) {
	var formData requests.SliderFormData

	formData.Nama = strings.TrimSpace(c.FormValue("name"))
	if formData.Nama == "" {
		return formData, slider_errors.ErrApiInvalidName(c)
	}

	file, err := c.FormFile("image_slider")
	if err != nil {
		if requireImage {
			h.logger.Debug("Image upload error", zap.Error(err))
			return formData, slider_errors.ErrApiImageRequired(c)
		}
		return formData, nil
	}

	imagePath, err := h.upload_image.ProcessImageUpload(c, file)
	if err != nil {
		return formData, err
	}

	formData.FilePath = imagePath
	return formData, nil
}

func (s *sliderHandleApi) startTracingAndLogging(
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

func (s *sliderHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
