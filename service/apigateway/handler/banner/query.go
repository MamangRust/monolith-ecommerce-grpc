package bannerhandler

import (
	"fmt"
	"net/http"
	"strconv"

	banner_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/banner"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/banner"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)


type bannerQueryHandlerApi struct {
	client        pb.BannerQueryServiceClient
	logger        logger.LoggerInterface
	mapper        apimapper.BannerQueryResponseMapper
	cache         banner_cache.BannerQueryCache
	observability observability.TraceLoggerObservability
}

type bannerQueryHandleDeps struct {
	client        pb.BannerQueryServiceClient
	router        *echo.Echo
	logger        logger.LoggerInterface
	mapper        apimapper.BannerQueryResponseMapper
	cache         banner_cache.BannerQueryCache
	observability observability.TraceLoggerObservability
}

func NewBannerQueryHandleApi(params *bannerQueryHandleDeps) *bannerQueryHandlerApi {
	handler := &bannerQueryHandlerApi{
		client:        params.client,
		logger:        params.logger,
		mapper:        params.mapper,
		cache:         params.cache,
		observability: params.observability,
	}

	routerBanner := params.router.Group("/api/banner-query")
	routerBanner.GET("", handler.FindAll)
	routerBanner.GET("/:id", handler.FindById)
	routerBanner.GET("/active", handler.FindByActive)
	routerBanner.GET("/trashed", handler.FindByTrashed)

	return handler
}


// @Security Bearer
// @Summary Find all banners
// @Tags Banner Query
// @Description Retrieve a list of all banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBanner "List of banners"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner-query [get]
func (h *bannerQueryHandlerApi) FindAll(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"FindAll",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	req := &requests.FindAllBanner{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedBanners(ctx, req); found {
		logSuccess("Serving from cache")
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllBannerRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindAll")
	}

	apiResponse := h.mapper.ToApiResponsePaginationBanner(res)
	h.cache.SetCachedBanners(ctx, req, apiResponse)

	logSuccess("Request completed successfully")
	return c.JSON(http.StatusOK, apiResponse)
}

func (h *bannerQueryHandlerApi) FindById(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"FindById",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		status = "error"
		return h.handleError(c, errors.NewBadRequestError("Invalid ID"), span, "FindById")
	}

	if cachedData, found := h.cache.GetCachedBanner(ctx, id); found {
		logSuccess("Serving from cache")
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdBannerRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindById")
	}

	apiResponse := h.mapper.ToApiResponseBanner(res)
	h.cache.SetCachedBanner(ctx, apiResponse)

	logSuccess("Request completed successfully")
	return c.JSON(http.StatusOK, apiResponse)
}

func (h *bannerQueryHandlerApi) FindByActive(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"FindByActive",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	req := &requests.FindAllBanner{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedActiveBanners(ctx, req); found {
		logSuccess("Serving from cache")
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllBannerRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByActive")
	}

	apiResponse := h.mapper.ToApiResponsePaginationBannerDeleteAt(res)
	h.cache.SetCachedActiveBanners(ctx, req, apiResponse)

	logSuccess("Request completed successfully")
	return c.JSON(http.StatusOK, apiResponse)
}

func (h *bannerQueryHandlerApi) FindByTrashed(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"FindByTrashed",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	req := &requests.FindAllBanner{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedTrashedBanners(ctx, req); found {
		logSuccess("Serving from cache")
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllBannerRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByTrashed")
	}

	apiResponse := h.mapper.ToApiResponsePaginationBannerDeleteAt(res)
	h.cache.SetCachedTrashedBanners(ctx, req, apiResponse)

	logSuccess("Request completed successfully")
	return c.JSON(http.StatusOK, apiResponse)
}

func (h *bannerQueryHandlerApi) handleError(c echo.Context, err error, span trace.Span, method string) error {
	appErr := errors.ParseGrpcError(err)
	traceID := span.SpanContext().TraceID().String()

	h.logger.Error(
		fmt.Sprintf("Banner query error in %s", method),
		zap.Error(err),
		zap.String("trace.id", traceID),
	)

	return errors.HandleApiError(c, appErr, traceID)
}

