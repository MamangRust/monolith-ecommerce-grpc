package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
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

type reviewHandleApi struct {
	client          pb.ReviewServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.ReviewResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerReview(
	router *echo.Echo,
	client pb.ReviewServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ReviewResponseMapper,
) *reviewHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "review_handler_requests_total",
			Help: "Total number of review requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "review_handler_request_duration_seconds",
			Help:    "Duration of review requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	reviewHandler := &reviewHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("review-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerreview := router.Group("/api/review")

	routerreview.GET("", reviewHandler.FindAll)
	routerreview.GET("/product/:id", reviewHandler.FindByProduct)
	routerreview.GET("/active", reviewHandler.FindByActive)
	routerreview.GET("/trashed", reviewHandler.FindByTrashed)

	routerreview.POST("/create", reviewHandler.Create)
	routerreview.POST("/update/:id", reviewHandler.Update)

	routerreview.POST("/trashed/:id", reviewHandler.TrashedReview)
	routerreview.POST("/restore/:id", reviewHandler.RestoreReview)
	routerreview.DELETE("/permanent/:id", reviewHandler.DeleteReviewPermanent)

	routerreview.POST("/restore/all", reviewHandler.RestoreAllReview)
	routerreview.POST("/permanent/all", reviewHandler.DeleteAllReviewPermanent)

	return reviewHandler
}

// @Security Bearer
// @Summary Find all review
// @Tags Review
// @Description Retrieve a list of all review
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of review"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review [get]
func (h *reviewHandleApi) FindAll(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAll"
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

		logError("Failed to retrieve review data", err, zap.Error(err))

		return review_errors.ErrApiFailedFindAllReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReview(res)

	logSuccess("Successfully retrieve review data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find reviews by product ID
// @Tags Review
// @Description Retrieve a list of reviews for a specific product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of reviews for the product"
// @Failure 400 {object} response.ErrorResponse "Invalid product ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/product/{id} [get]
func (h *reviewHandleApi) FindByProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return review_errors.ErrApiReviewInvalidProductId(c)
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")
	ctx := c.Request().Context()

	req := &pb.FindAllReviewProductRequest{
		ProductId: int32(id),
		Page:      int32(page),
		PageSize:  int32(pageSize),
		Search:    search,
	}

	res, err := h.client.FindByProduct(ctx, req)
	if err != nil {
		h.logger.Error("Failed to fetch reviews product", zap.Error(err))
		return review_errors.ErrApiFailedFindProductReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewsDetail(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find reviews by merchant ID
// @Tags Review
// @Description Retrieve a list of reviews for a specific merchant
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of reviews for the merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/merchant/{id} [get]
func (h *reviewHandleApi) FindByMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return review_errors.ErrApiReviewInvalidMerchantId(c)
	}

	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByMerchant"
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

	req := &pb.FindAllReviewMerchantRequest{
		MerchantId: int32(id),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve review data", err, zap.Error(err))

		return review_errors.ErrApiFailedFindMerchantReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewsDetail(res)

	logSuccess("Successfully retrieve review data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active review
// @Tags Review
// @Description Retrieve a list of active review
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDeleteAt "List of active review"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/active [get]
func (h *reviewHandleApi) FindByActive(c echo.Context) error {
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

		logError("Failed to retrieve review data", err, zap.Error(err))

		return review_errors.ErrApiFailedFindActiveReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDeleteAt(res)

	logSuccess("Successfully retrieve review data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed review records.
// @Summary Retrieve trashed review
// @Tags Review
// @Description Retrieve a list of trashed review records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDeleteAt "List of trashed review data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/trashed [get]
func (h *reviewHandleApi) FindByTrashed(c echo.Context) error {
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

		logError("Failed to retrieve review data", err, zap.Error(err))

		return review_errors.ErrApiFailedFindTrashedReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDeleteAt(res)

	logSuccess("Successfully retrieve review data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new review without image upload.
// @Summary Create a new review
// @Tags Review
// @Description Create a new review with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateReviewRequest true "review details"
// @Success 201 {object} response.ApiResponseReview "Successfully created review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create review"
// @Router /api/review/create [post]
func (h *reviewHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	userID, ok := c.Get("user_id").(int)

	if !ok || userID <= 0 {
		h.logger.Debug("Invalid or missing user ID from JWT")

		return review_errors.ErrApiReviewInvalidId(c)
	}

	var req requests.CreateReviewRequest

	if err := c.Bind(&req); err != nil {
		logError("Invalid request format", err, zap.Error(err))
		return review_errors.ErrApiBindCreateReview(c)
	}

	if err := req.Validate(); err != nil {
		logError("Validation failed", err, zap.Error(err))
		return review_errors.ErrApiValidateCreateReview(c)
	}

	grpcReq := &pb.CreateReviewRequest{
		UserId:    int32(userID),
		ProductId: int32(req.ProductID),
		Comment:   req.Comment,
		Rating:    int32(req.Rating),
	}

	res, err := h.client.Create(ctx, grpcReq)

	if err != nil {
		logError("review creation failed", err,
			zap.Error(err),
			zap.Any("request", req),
		)

		return review_errors.ErrApiFailedCreateReview(c)
	}

	logSuccess("successfully create review", zap.Int("user_id", userID))

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// Update handles the update of an existing review.
// @Summary Update an existing review
// @Tags Review
// @Description Update an existing review record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Param request body requests.UpdateReviewRequest true "review update details"
// @Success 200 {object} response.ApiResponseReview "Successfully updated review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update review"
// @Router /api/review/update/{id} [post]
func (h *reviewHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid id parameter", err, zap.Error(err))
		return review_errors.ErrApiReviewInvalidId(c)
	}

	var req requests.UpdateReviewRequest

	if err := c.Bind(&req); err != nil {
		logError("invalid body", err, zap.Error(err))

		return review_errors.ErrApiBindUpdateReview(c)
	}

	grpcReq := &pb.UpdateReviewRequest{
		ReviewId: int32(idInt),
		Name:     req.Name,
		Comment:  req.Comment,
		Rating:   int32(req.Rating),
	}

	res, err := h.client.Update(ctx, grpcReq)

	if err != nil {
		logError("review update failed", err,
			zap.Error(err),
			zap.Any("request", req),
		)

		return review_errors.ErrApiFailedUpdateReview(c)
	}

	logSuccess("successfully update review")

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// TrashedReview retrieves a trashed review record by its ID.
// @Summary Retrieve a trashed review
// @Tags Review
// @Description Retrieve a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewDeleteAt "Successfully retrieved trashed review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed review"
// @Router /api/review/trashed/{id} [get]
func (h *reviewHandleApi) TrashedReview(c echo.Context) error {
	const method = "Trashed"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid review ID format", err, zap.Error(err))
		return review_errors.ErrApiReviewInvalidId(c)
	}

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedReview(ctx, req)

	if err != nil {
		logError("Failed to archive review", err, zap.Error(err))
		return review_errors.ErrApiFailedTrashedReview(c)
	}

	so := h.mapping.ToApiResponseReviewDeleteAt(res)

	logSuccess("successfully trash review")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreReview restores a review record from the trash by its ID.
// @Summary Restore a trashed review
// @Tags Review
// @Description Restore a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} response.ApiResponseReviewDeleteAt "Successfully restored review"
// @Failure 400 {object} response.ErrorResponse "Invalid review ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore review"
// @Router /api/review/restore/{id} [post]
func (h *reviewHandleApi) RestoreReview(c echo.Context) error {
	const method = "Restore"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid review ID format", err, zap.Error(err))
		return review_errors.ErrApiReviewInvalidId(c)
	}

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreReview(ctx, req)

	if err != nil {
		logError("Failed to restore review", err, zap.Error(err))
		return review_errors.ErrApiFailedRestoreReview(c)
	}

	so := h.mapping.ToApiResponseReviewDeleteAt(res)

	logSuccess("successfully restore review")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteReviewPermanent permanently deletes a review record by its ID.
// @Summary Permanently delete a review
// @Tags review
// @Description Permanently delete a review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewDelete "Successfully deleted review record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete review:"
// @Router /api/review/delete/{id} [delete]
func (h *reviewHandleApi) DeleteReviewPermanent(c echo.Context) error {
	const method = "DeleteReview"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid review ID format", err, zap.Error(err))
		return review_errors.ErrApiReviewInvalidId(c)
	}

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteReviewPermanent(ctx, req)

	if err != nil {
		logError("Failed to permanently delete review", err, zap.Error(err))
		return review_errors.ErrApiFailedDeleteReviewPermanent(c)
	}

	so := h.mapping.ToApiResponseReviewDelete(res)

	logSuccess("successfull delete review permanent")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllReview restores a review record from the trash by its ID.
// @Summary Restore a trashed review
// @Tags review
// @Description Restore a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewAll "Successfully restored review all"
// @Failure 400 {object} response.ErrorResponse "Invalid review ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore review"
// @Router /api/review/restore/all [post]
func (h *reviewHandleApi) RestoreAllReview(c echo.Context) error {
	const method = "RestoreAll"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllReview(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Bulk review restoration failed", err, zap.Error(err))
		return review_errors.ErrApiFailedRestoreAllReviews(c)
	}

	so := h.mapping.ToApiResponseReviewAll(res)

	logSuccess("Successfully restored all review")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllReviewPermanent permanently deletes a review record by its ID.
// @Summary Permanently delete a review
// @Tags Review
// @Description Permanently delete a review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewAll "Successfully deleted review record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete review:"
// @Router /api/review/delete/all [post]
func (h *reviewHandleApi) DeleteAllReviewPermanent(c echo.Context) error {
	const method = "DeleteAll"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
	)
	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllReviewPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Bulk review deletion failed", err, zap.Error(err))
		return review_errors.ErrApiFailedDeleteAllReviewsPermanent(c)
	}

	so := h.mapping.ToApiResponseReviewAll(res)

	logSuccess("Successfully deleted all review permanently")

	return c.JSON(http.StatusOK, so)
}

func (s *reviewHandleApi) startTracingAndLogging(
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

func (s *reviewHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
