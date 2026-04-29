package orderitemhandler

import (
	"net/http"
	"strconv"

	orderitem_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/order_item"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/order_item"
	"github.com/labstack/echo/v4"
)

type orderItemQueryHandlerApi struct {
	client pb.OrderItemQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.OrderItemQueryResponseMapper
	cache  orderitem_cache.OrderItemQueryCache
}

type orderItemQueryHandleDeps struct {
	client pb.OrderItemQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.OrderItemQueryResponseMapper
	cache  orderitem_cache.OrderItemQueryCache
}

func NewOrderItemQueryHandleApi(params *orderItemQueryHandleDeps) *orderItemQueryHandlerApi {
	handler := &orderItemQueryHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
		cache:  params.cache,
	}

	routerOrderItem := params.router.Group("/api/order-item")
	routerOrderItem.GET("", handler.FindAll)
	routerOrderItem.GET("/order/:order_id", handler.FindOrderItemByOrder)
	routerOrderItem.GET("/active", handler.FindByActive)
	routerOrderItem.GET("/trashed", handler.FindByTrashed)

	return handler
}

// @Security Bearer
// @Summary Find all order items
// @Tags Order Item Query
// @Description Retrieve a list of all order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItem "List of order items"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item [get]
func (h *orderItemQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllOrderItems{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedOrderItemsAll(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllOrderItemRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return h.handleGrpcError(err, "FindAll") }

	apiResponse := h.mapper.ToApiResponsePaginationOrderItem(res)
	h.cache.SetCachedOrderItemsAll(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active order items
// @Tags Order Item Query
// @Description Retrieve a list of active order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItemDeleteAt "List of active order items"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/active [get]
func (h *orderItemQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllOrderItems{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedOrderItemActive(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllOrderItemRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return h.handleGrpcError(err, "FindByActive") }

	apiResponse := h.mapper.ToApiResponsePaginationOrderItemDeleteAt(res)
	h.cache.SetCachedOrderItemActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed order items
// @Tags Order Item Query
// @Description Retrieve a list of trashed order item records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItemDeleteAt "List of trashed order item data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/trashed [get]
func (h *orderItemQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllOrderItems{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedOrderItemTrashed(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllOrderItemRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return h.handleGrpcError(err, "FindByTrashed") }

	apiResponse := h.mapper.ToApiResponsePaginationOrderItemDeleteAt(res)
	h.cache.SetCachedOrderItemTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find order items by Order ID
// @Tags Order Item Query
// @Description Retrieve all items belonging to a specific order
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.ApiResponsesOrderItem "List of items in order"
// @Failure 400 {object} errors.ErrorResponse "Invalid Order ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/order/{order_id} [get]
func (h *orderItemQueryHandlerApi) FindOrderItemByOrder(c echo.Context) error {
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil || orderId <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid Order ID") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedOrderItems(ctx, orderId); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{Id: int32(orderId)})
	if err != nil { return h.handleGrpcError(err, "FindOrderItemByOrder") }

	apiResponse := h.mapper.ToApiResponsesOrderItem(res)
	h.cache.SetCachedOrderItems(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h *orderItemQueryHandlerApi) handleGrpcError(err error, operation string) error {
	h.logger.Error("Failed to " + operation)
	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to "+operation)
}
