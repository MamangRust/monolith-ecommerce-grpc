package reviewhandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	review_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/review"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/review"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type reviewQueryHandleApi struct {
	client        pb.ReviewQueryServiceClient
	logger        logger.LoggerInterface
	mapper        apimapper.ReviewQueryResponseMapper
	cache         review_cache.ReviewQueryCache
	observability observability.TraceLoggerObservability
}

type reviewQueryHandleDeps struct {
	client        pb.ReviewQueryServiceClient
	router        *echo.Echo
	logger        logger.LoggerInterface
	mapper        apimapper.ReviewQueryResponseMapper
	cache         review_cache.ReviewQueryCache
	observability observability.TraceLoggerObservability
}

func NewReviewQueryHandleApi(deps *reviewQueryHandleDeps) {
	handler := &reviewQueryHandleApi{
		client:        deps.client,
		logger:        deps.logger,
		mapper:        deps.mapper,
		cache:         deps.cache,
		observability: deps.observability,
	}

	router := deps.router.Group("/api/review-query")
	router.GET("", handler.FindAll)
	router.GET("/product/:id", handler.FindByProduct)
	router.GET("/merchant/:id", handler.FindByMerchant)
	router.GET("/active", handler.FindByActive)
	router.GET("/trashed", handler.FindByTrashed)
}

// @Security Bearer
// @Summary Find all reviews
// @Tags Review Query
// @Description Retrieve a list of all reviews
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of reviews"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review data"
// @Router /api/review-query [get]
func (h *reviewQueryHandleApi) FindAll(c echo.Context) error {
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
	if cached, found := h.cache.GetReviewAllCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindAllReviews")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindAll")
	}

	response := h.mapper.ToApiResponsePaginationReview(res)
	h.cache.SetReviewAllCache(ctx, req, response)

	logSuccess("Successfully retrieved reviews")
	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Find reviews by Product ID
// @Tags Review Query
// @Description Retrieve reviews for a specific product including review details
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewsDetail "List of reviews with details"
// @Failure 400 {object} errors.ErrorResponse "Invalid Product ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review data"
// @Router /api/review-query/product/{id} [get]
func (h *reviewQueryHandleApi) FindByProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Product ID")
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

	grpcReq := &pb.FindAllReviewProductRequest{
		ProductId: int32(id),
		Page:      int32(page),
		PageSize:  int32(pageSize),
		Search:    search,
	}

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByProductReviews")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindByProduct(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByProduct")
	}

	response := h.mapper.ToApiResponsePaginationReviewsDetail(res)

	logSuccess("Successfully retrieved reviews for product")
	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Find reviews by Merchant ID
// @Tags Review Query
// @Description Retrieve reviews for a specific merchant including review details
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewsDetail "List of reviews with details"
// @Failure 400 {object} errors.ErrorResponse "Invalid Merchant ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review data"
// @Router /api/review-query/merchant/{id} [get]
func (h *reviewQueryHandleApi) FindByMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Merchant ID")
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

	req := &requests.FindAllReviewByMerchant{
		MerchantID: id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetReviewByMerchantCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllReviewMerchantRequest{
		MerchantId: int32(id),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByMerchantReviews")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindByMerchant(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByMerchant")
	}

	response := h.mapper.ToApiResponsePaginationReviewsDetail(res)
	h.cache.SetReviewByMerchantCache(ctx, req, response)

	logSuccess("Successfully retrieved reviews for merchant")
	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Retrieve active reviews
// @Tags Review Query
// @Description Retrieve a list of active reviews
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDeleteAt "List of active reviews"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review data"
// @Router /api/review-query/active [get]
func (h *reviewQueryHandleApi) FindByActive(c echo.Context) error {
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
	if cached, found := h.cache.GetReviewActiveCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByActiveReviews")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByActive")
	}

	response := h.mapper.ToApiResponsePaginationReviewDeleteAt(res)
	h.cache.SetReviewActiveCache(ctx, req, response)

	logSuccess("Successfully retrieved active reviews")
	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Retrieve trashed reviews
// @Tags Review Query
// @Description Retrieve a list of trashed review records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDeleteAt "List of trashed reviews"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve review data"
// @Router /api/review-query/trashed [get]
func (h *reviewQueryHandleApi) FindByTrashed(c echo.Context) error {
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
	if cached, found := h.cache.GetReviewTrashedCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByTrashedReviews")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByTrashed")
	}

	response := h.mapper.ToApiResponsePaginationReviewDeleteAt(res)
	h.cache.SetReviewTrashedCache(ctx, req, response)

	logSuccess("Successfully retrieved trashed reviews")
	return c.JSON(http.StatusOK, response)
}

func (h *reviewQueryHandleApi) handleError(c echo.Context, err error, span trace.Span, method string) error {
	appErr := errors.ParseGrpcError(err)
	traceID := span.SpanContext().TraceID().String()

	h.logger.Error(
		"Review query error in "+method,
		zap.Error(err),
		zap.String("trace.id", traceID),
	)

	return errors.HandleApiError(c, appErr, traceID)
}
