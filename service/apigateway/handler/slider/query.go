package sliderhandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	slider_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/slider"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/slider"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type sliderQueryHandleApi struct {
	client pb.SliderQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.SliderQueryResponseMapper
	cache  slider_cache.SliderQueryCache
}

type sliderQueryHandleDeps struct {
	client pb.SliderQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.SliderQueryResponseMapper
	cache  slider_cache.SliderQueryCache
}

func NewSliderQueryHandleApi(deps *sliderQueryHandleDeps) {
	handler := &sliderQueryHandleApi{
		client: deps.client,
		logger: deps.logger,
		mapper: deps.mapper,
		cache:  deps.cache,
	}

	router := deps.router.Group("/api/slider-query")
	router.GET("", handler.FindAll)
	router.GET("/:id", handler.FindById)
	router.GET("/active", handler.FindByActive)
	router.GET("/trashed", handler.FindByTrashed)
}

// @Security Bearer
// @Summary Find all sliders
// @Tags Slider Query
// @Description Retrieve a list of all sliders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSlider "List of sliders"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider-query [get]
func (h *sliderQueryHandleApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetSliderAllCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAll")
	}

	response := h.mapper.ToApiResponsePaginationSlider(res)
	h.cache.SetSliderAllCache(ctx, req, response)

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Find slider by ID
// @Tags Slider Query
// @Description Retrieve a slider by ID
// @Accept json
// @Produce json
// @Param id path int true "Slider ID"
// @Success 200 {object} response.ApiResponseSlider "Slider data"
// @Failure 400 {object} errors.ErrorResponse "Invalid slider ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider-query/{id} [get]
func (h *sliderQueryHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetCachedSliderCache(ctx, id); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdSliderRequest{Id: int32(id)})
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	response := h.mapper.ToApiResponseSlider(res)
	h.cache.SetCachedSliderCache(ctx, response)

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Retrieve active sliders
// @Tags Slider Query
// @Description Retrieve a list of active sliders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSliderDeleteAt "List of active sliders"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider-query/active [get]
func (h *sliderQueryHandleApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetSliderActiveCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	response := h.mapper.ToApiResponsePaginationSliderDeleteAt(res)
	h.cache.SetSliderActiveCache(ctx, req, response)

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Retrieve trashed sliders
// @Tags Slider Query
// @Description Retrieve a list of trashed slider records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSliderDeleteAt "List of trashed sliders"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve slider data"
// @Router /api/slider-query/trashed [get]
func (h *sliderQueryHandleApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetSliderTrashedCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllSliderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	response := h.mapper.ToApiResponsePaginationSliderDeleteAt(res)
	h.cache.SetSliderTrashedCache(ctx, req, response)

	return c.JSON(http.StatusOK, response)
}

func (h *sliderQueryHandleApi) handleGrpcError(err error, operation string) error {
	h.logger.Error("Failed to "+operation, zap.Error(err))
	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to "+operation)
}
