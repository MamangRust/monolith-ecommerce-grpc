package merchantawardhandler

import (
	"net/http"
	"strconv"

	merchantawards_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_awards"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_award"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/labstack/echo/v4"
)

type merchantAwardQueryHandlerApi struct {
	client pb.MerchantAwardQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantAwardQueryResponseMapper
	cache  merchantawards_cache.MerchantAwardQueryCache
}

type merchantAwardQueryHandleDeps struct {
	client pb.MerchantAwardQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.MerchantAwardQueryResponseMapper
	cache  merchantawards_cache.MerchantAwardQueryCache
}

func NewMerchantAwardQueryHandleApi(params *merchantAwardQueryHandleDeps) *merchantAwardQueryHandlerApi {
	handler := &merchantAwardQueryHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
		cache:  params.cache,
	}

	routerAward := params.router.Group("/api/merchant-award-query")
	routerAward.GET("", handler.FindAll)
	routerAward.GET("/:id", handler.FindById)
	routerAward.GET("/active", handler.FindByActive)
	routerAward.GET("/trashed", handler.FindByTrashed)

	return handler
}

// @Security Bearer
// @Summary Find all merchant awards
// @Tags Merchant Award Query
// @Description Retrieve a list of all merchant awards/certifications
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantAward "List of merchant awards"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant award data"
// @Router /api/merchant-award-query [get]
func (h *merchantAwardQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantAwardAll(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantAward(res)
	h.cache.SetCachedMerchantAwardAll(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find merchant award by ID
// @Tags Merchant Award Query
// @Description Retrieve a merchant award by ID
// @Accept json
// @Produce json
// @Param id path int true "Award ID"
// @Success 200 {object} response.ApiResponseMerchantAward "Merchant award data"
// @Failure 400 {object} errors.ErrorResponse "Invalid award ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant award data"
// @Router /api/merchant-award-query/{id} [get]
func (h *merchantAwardQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedMerchantAward(ctx, id); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdMerchantAwardRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseMerchantAward(res)
	h.cache.SetCachedMerchantAward(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active merchant awards
// @Tags Merchant Award Query
// @Description Retrieve a list of active merchant awards
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantAwardDeleteAt "List of active merchant awards"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant award data"
// @Router /api/merchant-award-query/active [get]
func (h *merchantAwardQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantAwardActive(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantAwardDeleteAt(res)
	h.cache.SetCachedMerchantAwardActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed merchant awards
// @Tags Merchant Award Query
// @Description Retrieve a list of trashed merchant award records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantAwardDeleteAt "List of trashed merchant award data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant award data"
// @Router /api/merchant-award-query/trashed [get]
func (h *merchantAwardQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantAwardTrashed(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantAwardDeleteAt(res)
	h.cache.SetCachedMerchantAwardTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

