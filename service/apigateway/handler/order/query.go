package orderhandler

import (
	"net/http"
	"strconv"

	order_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/order"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/order"
	"github.com/labstack/echo/v4"
)

type orderQueryHandlerApi struct {
	client pb.OrderQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.OrderQueryResponseMapper
	cache  order_cache.OrderQueryCache
}

type orderQueryHandleDeps struct {
	client pb.OrderQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.OrderQueryResponseMapper
	cache  order_cache.OrderQueryCache
}

func NewOrderQueryHandleApi(params *orderQueryHandleDeps) *orderQueryHandlerApi {
	handler := &orderQueryHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
		cache:  params.cache,
	}

	routerOrder := params.router.Group("/api/order-query")
	routerOrder.GET("", handler.FindAll)
	routerOrder.GET("/:id", handler.FindById)
	routerOrder.GET("/active", handler.FindByActive)
	routerOrder.GET("/trashed", handler.FindByTrashed)

	return handler
}

// @Security Bearer
// @Summary Find all orders
// @Tags Order Query
// @Description Retrieve a list of all orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrder "List of orders"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve order data"
// @Router /api/order-query [get]
func (h *orderQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllOrder{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetOrderAllCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllOrderRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationOrder(res)
	h.cache.SetOrderAllCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find order by ID
// @Tags Order Query
// @Description Retrieve an order by ID
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrder "Order data"
// @Failure 400 {object} errors.ErrorResponse "Invalid order ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve order data"
// @Router /api/order-query/{id} [get]
func (h *orderQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedOrderCache(ctx, id); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdOrderRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseOrder(res)
	h.cache.SetCachedOrderCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active orders
// @Tags Order Query
// @Description Retrieve a list of active orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderDeleteAt "List of active orders"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve order data"
// @Router /api/order-query/active [get]
func (h *orderQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllOrder{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetOrderActiveCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllOrderRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationOrderDeleteAt(res)
	h.cache.SetOrderActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed orders
// @Tags Order Query
// @Description Retrieve a list of trashed order records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderDeleteAt "List of trashed order data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve order data"
// @Router /api/order-query/trashed [get]
func (h *orderQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllOrder{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetOrderTrashedCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllOrderRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationOrderDeleteAt(res)
	h.cache.SetOrderTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

