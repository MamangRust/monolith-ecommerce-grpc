package merchanthandler

import (
	"net/http"
	"strconv"

	merchant_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	"github.com/labstack/echo/v4"
)




type merchantQueryHandlerApi struct {
	client pb.MerchantQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantQueryResponseMapper
	cache  merchant_cache.MerchantQueryCache
	apiHandler errors.ApiHandler
}

type merchantQueryHandleDeps struct {
	client     pb.MerchantQueryServiceClient
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     apimapper.MerchantQueryResponseMapper
	cache      merchant_cache.MerchantQueryCache
	apiHandler errors.ApiHandler
}

func NewMerchantQueryHandleApi(params *merchantQueryHandleDeps) *merchantQueryHandlerApi {
	merchantQueryHandler := &merchantQueryHandlerApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	routerMerchant := params.router.Group("/api/merchant-query")

	routerMerchant.GET("", merchantQueryHandler.FindAll)
	routerMerchant.GET("/:id", merchantQueryHandler.FindById)
	routerMerchant.GET("/active", merchantQueryHandler.FindByActive)
	routerMerchant.GET("/trashed", merchantQueryHandler.FindByTrashed)

	return merchantQueryHandler
}

// @Security Bearer
// @Summary Find all merchants
// @Tags Merchant Query
// @Description Retrieve a list of all merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchant "List of merchants"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query [get]
func (h *merchantQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchants(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationMerchant(res)
	h.cache.SetCachedMerchants(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags Merchant Query
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Merchant data"
// @Failure 400 {object} errors.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query/{id} [get]
func (h *merchantQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedMerchant(ctx, id); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdMerchantRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponseMerchant(res)
	h.cache.SetCachedMerchant(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active merchants
// @Tags Merchant Query
// @Description Retrieve a list of active merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDeleteAt "List of active merchants"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query/active [get]
func (h *merchantQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantActive(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	// Note: MerchantQueryResponseMapper might need ToApiResponsePaginationMerchantDeleteAt
	// consolidated into it, similar to Category.
	apiResponse := h.mapper.ToApiResponsePaginationMerchantDeleteAt(res)
	h.cache.SetCachedMerchantActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed merchants
// @Tags Merchant Query
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDeleteAt "List of trashed merchant data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query/trashed [get]
func (h *merchantQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedMerchantTrashed(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationMerchantDeleteAt(res)
	h.cache.SetCachedMerchantTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}


