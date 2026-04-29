package merchantbusinesshandler

import (
	"net/http"
	"strconv"

	merchantbusiness_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_business"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_business"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/labstack/echo/v4"
)

type merchantBusinessQueryHandlerApi struct {
	client pb.MerchantBusinessQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantBusinessQueryResponseMapper
	cache  merchantbusiness_cache.MerchantBusinessQueryCache
}

type merchantBusinessQueryHandleDeps struct {
	client pb.MerchantBusinessQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.MerchantBusinessQueryResponseMapper
	cache  merchantbusiness_cache.MerchantBusinessQueryCache
}

func NewMerchantBusinessQueryHandleApi(params *merchantBusinessQueryHandleDeps) *merchantBusinessQueryHandlerApi {
	handler := &merchantBusinessQueryHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
		cache:  params.cache,
	}

	routerBusiness := params.router.Group("/api/merchant-business-query")
	routerBusiness.GET("", handler.FindAll)
	routerBusiness.GET("/:id", handler.FindById)
	routerBusiness.GET("/active", handler.FindByActive)
	routerBusiness.GET("/trashed", handler.FindByTrashed)

	return handler
}

// @Security Bearer
// @Summary Find all merchant business details
// @Tags Merchant Business Query
// @Description Retrieve a list of all merchant business details
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusiness "List of merchant business details"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant business data"
// @Router /api/merchant-business-query [get]
func (h *merchantBusinessQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantBusinessAll(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantBusiness(res)
	h.cache.SetCachedMerchantBusinessAll(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find merchant business by ID
// @Tags Merchant Business Query
// @Description Retrieve a merchant business detail by ID
// @Accept json
// @Produce json
// @Param id path int true "Business ID"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Merchant business data"
// @Failure 400 {object} errors.ErrorResponse "Invalid business ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant business data"
// @Router /api/merchant-business-query/{id} [get]
func (h *merchantBusinessQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedMerchantBusiness(ctx, id); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdMerchantBusinessRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseMerchantBusiness(res)
	h.cache.SetCachedMerchantBusiness(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active merchant business details
// @Tags Merchant Business Query
// @Description Retrieve a list of active merchant business details
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusinessDeleteAt "List of active merchant business details"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant business data"
// @Router /api/merchant-business-query/active [get]
func (h *merchantBusinessQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantBusinessActive(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantBusinessDeleteAt(res)
	h.cache.SetCachedMerchantBusinessActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed merchant business details
// @Tags Merchant Business Query
// @Description Retrieve a list of trashed merchant business records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusinessDeleteAt "List of trashed merchant business data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant business data"
// @Router /api/merchant-business-query/trashed [get]
func (h *merchantBusinessQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantBusinessTrashed(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantBusinessDeleteAt(res)
	h.cache.SetCachedMerchantBusinessTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

