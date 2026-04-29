package reviewdetailhandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	reviewdetail_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/review_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/review_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewDetailQueryHandleApi struct {
	client        pb.ReviewDetailQueryServiceClient
	logger        logger.LoggerInterface
	mapper        apimapper.ReviewDetailQueryResponseMapper
	cache         reviewdetail_cache.ReviewDetailQueryCache
	observability observability.TraceLoggerObservability
}

type reviewDetailQueryHandleDeps struct {
	client        pb.ReviewDetailQueryServiceClient
	router        *echo.Echo
	logger        logger.LoggerInterface
	mapper        apimapper.ReviewDetailQueryResponseMapper
	cache         reviewdetail_cache.ReviewDetailQueryCache
	observability observability.TraceLoggerObservability
}

func NewReviewDetailQueryHandleApi(deps *reviewDetailQueryHandleDeps) {
	handler := &reviewDetailQueryHandleApi{
		client:        deps.client,
		logger:        deps.logger,
		mapper:        deps.mapper,
		cache:         deps.cache,
		observability: deps.observability,
	}

	router := deps.router.Group("/api/review-detail-query")
	router.GET("", handler.FindAll)
	router.GET("/:id", handler.FindById)
	router.GET("/active", handler.FindByActive)
	router.GET("/trashed", handler.FindByTrashed)
}

// @Security Bearer
// @Summary Find all review details
// @Tags Review Detail Query
// @Description Retrieve a list of all review details
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetails "List of review details"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review detail data"
// @Router /api/review-detail-query [get]
func (h *reviewDetailQueryHandleApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetReviewDetailAllCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindAllReviewDetails")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindAll")
	}

	response := h.mapper.ToApiResponsePaginationReviewDetail(res)
	h.cache.SetReviewDetailAllCache(ctx, req, response)

	logSuccess("Successfully retrieved review details")
	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Find review detail by ID
// @Tags Review Detail Query
// @Description Retrieve a review detail by ID
// @Accept json
// @Produce json
// @Param id path int true "Review Detail ID"
// @Success 200 {object} response.ApiResponseReviewDetail "Review detail data"
// @Failure 400 {object} errors.ErrorResponse "Invalid review detail ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review detail data"
// @Router /api/review-detail-query/{id} [get]
func (h *reviewDetailQueryHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetCachedReviewDetailCache(ctx, id); found {
		return c.JSON(http.StatusOK, cached)
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByIdReviewDetail")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindById(ctx, &pb.FindByIdReviewDetailRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindById")
	}

	response := h.mapper.ToApiResponseReviewDetail(res)
	h.cache.SetCachedReviewDetailCache(ctx, response)

	logSuccess("Successfully retrieved review detail by id")
	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Retrieve active review details
// @Tags Review Detail Query
// @Description Retrieve a list of active review details
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetailsDeleteAt "List of active review details"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review detail data"
// @Router /api/review-detail-query/active [get]
func (h *reviewDetailQueryHandleApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetReviewDetailActiveCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByActiveReviewDetails")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByActive")
	}

	response := h.mapper.ToApiResponsePaginationReviewDetailDeleteAt(res)
	h.cache.SetReviewDetailActiveCache(ctx, req, response)

	logSuccess("Successfully retrieved active review details")
	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Retrieve trashed review details
// @Tags Review Detail Query
// @Description Retrieve a list of trashed review detail records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetailsDeleteAt "List of trashed review detail data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review detail data"
// @Router /api/review-detail-query/trashed [get]
func (h *reviewDetailQueryHandleApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetReviewDetailTrashedCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByTrashedReviewDetails")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByTrashed")
	}

	response := h.mapper.ToApiResponsePaginationReviewDetailDeleteAt(res)
	h.cache.SetReviewDetailTrashedCache(ctx, req, response)

	logSuccess("Successfully retrieved trashed review details")
	return c.JSON(http.StatusOK, response)
}

func (h *reviewDetailQueryHandleApi) handleError(c echo.Context, err error, span trace.Span, method string) error {
	appErr := errors.ParseGrpcError(err)
	traceID := span.SpanContext().TraceID().String()

	h.logger.Error(
		"Review detail query error in "+method,
		zap.Error(err),
		zap.String("trace.id", traceID),
	)

	return errors.HandleApiError(c, appErr, traceID)
}
