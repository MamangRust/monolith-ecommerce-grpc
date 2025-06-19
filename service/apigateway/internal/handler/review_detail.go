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
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
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

type reviewDetailHandleApi struct {
	client          pb.ReviewDetailServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.ReviewDetailResponseMapper
	mappingReview   response_api.ReviewResponseMapper
	upload_image    upload_image.ImageUploads
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerReviewDetail(
	router *echo.Echo,
	client pb.ReviewDetailServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ReviewDetailResponseMapper,
	mappingReview response_api.ReviewResponseMapper,
	upload_image upload_image.ImageUploads,
) *reviewDetailHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_detail_handler_requests_total",
			Help: "Total number of review_detail requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "review_detail_handler_request_duration_seconds",
			Help:    "Duration of review_detail requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	reviewDetailHandler := &reviewDetailHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingReview:   mappingReview,
		upload_image:    upload_image,
		trace:           otel.Tracer("review_detail-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/review-detail")

	routercategory.GET("", reviewDetailHandler.FindAllReviewDetail)
	routercategory.GET("/:id", reviewDetailHandler.FindById)
	routercategory.GET("/active", reviewDetailHandler.FindByActive)
	routercategory.GET("/trashed", reviewDetailHandler.FindByTrashed)

	routercategory.POST("/create", reviewDetailHandler.Create)
	routercategory.POST("/update/:id", reviewDetailHandler.Update)

	routercategory.POST("/trashed/:id", reviewDetailHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", reviewDetailHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", reviewDetailHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", reviewDetailHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", reviewDetailHandler.DeleteAllMerchantPermanent)

	return reviewDetailHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags ReviewDetail
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetails "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail [get]
func (h *reviewDetailHandleApi) FindAllReviewDetail(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllReviewDetail"
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

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "error"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedFindAllReviewDetails(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDetail(res)

	logSuccess("Successfully retrieved merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags ReviewDetail
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetail "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/{id} [get]
func (h *reviewDetailHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID", err, zap.Error(err))

		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return reviewdetail_errors.ErrApiReviewDetailNotFound(c)
	}

	so := h.mapping.ToApiResponseReviewDetail(res)

	logSuccess("Successfully retrieved merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags ReviewDetail
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetailsDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/active [get]
func (h *reviewDetailHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "error"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedFindActiveReviewDetails(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDetailDeleteAt(res)

	logSuccess("Successfully retrieved merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags ReviewDetail
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationReviewDetailsDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/trashed [get]
func (h *reviewDetailHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedFindTrashedReviewDetails(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDetailDeleteAt(res)

	logSuccess("Successfully retrieved merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new merchant review detail.
// @Summary Create a new merchant review detail
// @Tags ReviewDetail
// @Description Create a new merchant review detail with the provided details
// @Accept multipart/form-data
// @Produce json
// @Param type formData string true "Type"
// @Param url formData file true "url"
// @Param caption formData string true "Product name"
// @Success 200 {object} response.ApiResponseReviewDetail "Successfully created review detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create review detail"
// @Router /api/review-detail/create [post]
func (h *reviewDetailHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	formData, err := h.parseReviewDetailForm(c)

	if err != nil {
		status = "error"
		logError("Review detail creation failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiInvalidBody(c)
	}

	req := &pb.CreateReviewDetailRequest{
		ReviewId: int32(formData.ReviewID),
		Type:     formData.Type,
		Url:      formData.Url,
		Caption:  formData.Caption,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		status = "error"

		logError("Review detail creation failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedCreateReviewDetail(c)
	}

	so := h.mapping.ToApiResponseReviewDetail(res)

	logSuccess("Successfully created review detail", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing merchant review detail.
// @Summary Update an existing merchant review detail
// @Tags ReviewDetail
// @Description Update an existing merchant review detail with the provided details
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Review Detail ID"
// @Param type formData string true "Type"
// @Param url formData file true "url"
// @Param caption formData string true "Product name"
// @Param request body requests.UpdateReviewDetailRequest true "Update review detail request"
// @Success 200 {object} response.ApiResponseReviewDetail "Successfully updated review detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update review detail"
// @Router /api/review-detail/update/{id} [post]
func (h *reviewDetailHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"
		logError("Review detail update failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	formData, err := h.parseReviewDetailForm(c)
	if err != nil {
		status = "error"

		logError("Review detail update failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiInvalidBody(c)
	}

	req := &pb.UpdateReviewDetailRequest{
		ReviewDetailId: int32(idInt),
		Type:           formData.Type,
		Url:            formData.Url,
		Caption:        formData.Caption,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		status = "error"
		logError("Review detail update failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedUpdateReviewDetail(c)
	}

	so := h.mapping.ToApiResponseReviewDetail(res)

	logSuccess("Successfully updated review detail", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags ReviewDetail
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetailDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/review-detail/trashed/{id} [get]
func (h *reviewDetailHandleApi) TrashedMerchant(c echo.Context) error {
	const method = "TrashedMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Review detail trashed failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedReviewDetail(ctx, req)

	if err != nil {
		status = "error"

		logError("Review detail trashed failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedTrashedReviewDetail(c)
	}

	so := h.mapping.ToApiResponseReviewDetailDeleteAt(res)

	logSuccess("Successfully trashed review detail", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags ReviewDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetailDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/review-detail/restore/{id} [post]
func (h *reviewDetailHandleApi) RestoreMerchant(c echo.Context) error {
	const method = "RestoreMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Review detail restore failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreReviewDetail(ctx, req)

	if err != nil {
		status = "error"

		logError("Review detail restore failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedRestoreReviewDetail(c)
	}

	so := h.mapping.ToApiResponseReviewDetailDeleteAt(res)

	logSuccess("Successfully restored review detail", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags ReviewDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/review-detail/delete/{id} [delete]
func (h *reviewDetailHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	const method = "DeleteMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Review detail permanent delete failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteReviewDetailPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Review detail permanent delete failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedDeleteReviewDetailPermanent(c)
	}

	so := h.mappingReview.ToApiResponseReviewDelete(res)

	logSuccess("Successfully deleted review detail permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags ReviewDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/review-detail/restore/all [post]
func (h *reviewDetailHandleApi) RestoreAllMerchant(c echo.Context) error {
	const method = "RestoreAllMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllReviewDetail(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Review detail restore all failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedRestoreAllReviewDetail(c)
	}

	so := h.mappingReview.ToApiResponseReviewAll(res)

	logSuccess("Successfully restored all review detail", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags ReviewDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/review-detail/delete/all [post]
func (h *reviewDetailHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	const method = "DeleteAllMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllReviewDetailPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Review detail permanent delete all failed", err, zap.Error(err))

		return reviewdetail_errors.ErrApiFailedDeleteAllReviewDetailPermanent(c)
	}

	so := h.mappingReview.ToApiResponseReviewAll(res)

	logSuccess("Successfully deleted all review detail permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (h *reviewDetailHandleApi) parseReviewDetailForm(c echo.Context) (requests.ReviewDetailFormData, error) {
	var formData requests.ReviewDetailFormData
	var err error

	formData.ReviewID, err = strconv.Atoi(c.FormValue("review_id"))
	if err != nil || formData.ReviewID <= 0 {
		return formData, reviewdetail_errors.ErrApiInvalidReviewId(c)
	}

	formData.Type = strings.TrimSpace(c.FormValue("type"))
	if formData.Type == "" {
		return formData, reviewdetail_errors.ErrApiReviewDetailTypeRequired(c)
	}

	formData.Caption = strings.TrimSpace(c.FormValue("caption"))
	if formData.Caption == "" {
		return formData, reviewdetail_errors.ErrApiReviewDetailCaptionRequired(c)
	}

	file, err := c.FormFile("url")
	if err != nil {
		return formData, reviewdetail_errors.ErrApiReviewDetailFileRequired(c)
	}

	uploadPath, err := h.upload_image.ProcessImageUpload(c, file)
	if err != nil {
		return formData, err
	}
	formData.Url = uploadPath

	return formData, nil
}

func (s *reviewDetailHandleApi) startTracingAndLogging(
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
		s.logger.Info(msg, fields...)
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

func (s *reviewDetailHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
