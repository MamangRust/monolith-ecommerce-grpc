package merchantdetailhandler

import (
	"net/http"
	"strconv"

	merchant_detail_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_detail"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_detail"
	"github.com/labstack/echo/v4"
)

type merchantDetailQueryHandlerApi struct {
	client     pb.MerchantDetailQueryServiceClient
	logger     logger.LoggerInterface
	mapper     apimapper.MerchantDetailQueryResponseMapper
	cache      merchant_detail_cache.MerchantDetailQueryCache
	apiHandler sharedErrors.ApiHandler
}

type merchantDetailQueryHandleDeps struct {
	client     pb.MerchantDetailQueryServiceClient
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     apimapper.MerchantDetailQueryResponseMapper
	cache      merchant_detail_cache.MerchantDetailQueryCache
	apiHandler sharedErrors.ApiHandler
}

func NewMerchantDetailQueryHandleApi(params *merchantDetailQueryHandleDeps) *merchantDetailQueryHandlerApi {
	handler := &merchantDetailQueryHandlerApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	router := params.router.Group("/api/merchant-detail-query")

	router.GET("", handler.FindAll)
	router.GET("/:id", handler.FindById)
	router.GET("/active", handler.FindByActive)
	router.GET("/trashed", handler.FindByTrashed)

	return handler
}

// @Security Bearer
// @Summary Find all merchant details
// @Tags Merchant Detail Query
// @Description Retrieve a list of all merchant details
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDetail "List of merchant details"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant detail data"
// @Router /api/merchant-detail-query [get]
func (h *merchantDetailQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	cacheReq := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantDetailAll(ctx, cacheReq); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllMerchantRequest{Page: int32(page), PageSize: int32(pageSize), Search: search})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantDetail(res)
	h.cache.SetCachedMerchantDetailAll(ctx, cacheReq, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find merchant detail by ID
// @Tags Merchant Detail Query
// @Description Retrieve a merchant detail by ID with all relations
// @Accept json
// @Produce json
// @Param id path int true "Detail ID"
// @Success 200 {object} response.ApiResponseMerchantDetailRelation "Merchant detail data"
// @Failure 400 {object} errors.ErrorResponse "Invalid detail ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant detail data"
// @Router /api/merchant-detail-query/{id} [get]
func (h *merchantDetailQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return sharedErrors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedMerchantDetailRelation(ctx, id); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdMerchantDetailRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseMerchantDetailRelation(res)
	h.cache.SetCachedMerchantDetailRelation(ctx, id, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active merchant details
// @Tags Merchant Detail Query
// @Description Retrieve a list of active merchant details
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDetailDeleteAt "List of active merchant details"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant detail data"
// @Router /api/merchant-detail-query/active [get]
func (h *merchantDetailQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	cacheReq := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantDetailActive(ctx, cacheReq); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllMerchantRequest{Page: int32(page), PageSize: int32(pageSize), Search: search})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantDetailDeleteAt(res)
	h.cache.SetCachedMerchantDetailActive(ctx, cacheReq, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed merchant details
// @Tags Merchant Detail Query
// @Description Retrieve a list of trashed merchant detail records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDetailDeleteAt "List of trashed merchant detail data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant detail data"
// @Router /api/merchant-detail-query/trashed [get]
func (h *merchantDetailQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	cacheReq := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantDetailTrashed(ctx, cacheReq); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllMerchantRequest{Page: int32(page), PageSize: int32(pageSize), Search: search})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantDetailDeleteAt(res)
	h.cache.SetCachedMerchantDetailTrashed(ctx, cacheReq, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

