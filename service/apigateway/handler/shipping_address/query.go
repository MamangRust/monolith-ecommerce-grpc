package shippingaddresshandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	shippingaddress_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/shipping_address"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/shipping_address"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type shippingAddressQueryHandleApi struct {
	client pb.ShippingQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.ShippingAddressQueryResponseMapper
	cache  shippingaddress_cache.ShippingAddressQueryCache
}

type shippingAddressQueryHandleDeps struct {
	client pb.ShippingQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.ShippingAddressQueryResponseMapper
	cache  shippingaddress_cache.ShippingAddressQueryCache
}

func NewShippingAddressQueryHandleApi(deps *shippingAddressQueryHandleDeps) {
	handler := &shippingAddressQueryHandleApi{
		client: deps.client,
		logger: deps.logger,
		mapper: deps.mapper,
		cache:  deps.cache,
	}

	router := deps.router.Group("/api/shipping-address-query")
	router.GET("", handler.FindAll)
	router.GET("/:id", handler.FindById)
	router.GET("/order/:id", handler.FindByOrder)
	router.GET("/active", handler.FindByActive)
	router.GET("/trashed", handler.FindByTrashed)
}

// @Security Bearer
// @Summary Find all shipping addresses
// @Tags Shipping Address Query
// @Description Retrieve a list of all shipping addresses
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddress "List of shipping addresses"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address-query [get]
func (h *shippingAddressQueryHandleApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetShippingAddressAllCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindAll")
	}

	response := h.mapper.ToApiResponsePaginationShippingAddress(res)
	h.cache.SetShippingAddressAllCache(ctx, req, response)

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Find shipping address by ID
// @Tags Shipping Address Query
// @Description Retrieve a shipping address by ID
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddress "Shipping address data"
// @Failure 400 {object} errors.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address-query/{id} [get]
func (h *shippingAddressQueryHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetCachedShippingAddressCache(ctx, id); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindByIdShippingRequest{Id: int32(id)}
	res, err := h.client.FindById(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindById")
	}

	response := h.mapper.ToApiResponseShippingAddress(res)
	h.cache.SetCachedShippingAddressCache(ctx, response)

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Find shipping address by Order ID
// @Tags Shipping Address Query
// @Description Retrieve a shipping address for a specific order
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseShippingAddress "Shipping address data"
// @Failure 400 {object} errors.ErrorResponse "Invalid order ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address-query/order/{id} [get]
func (h *shippingAddressQueryHandleApi) FindByOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetCachedShippingAddressByOrderCache(ctx, id); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindByIdShippingRequest{Id: int32(id)}
	res, err := h.client.FindByOrder(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByOrder")
	}

	response := h.mapper.ToApiResponseShippingAddress(res)
	h.cache.SetCachedShippingAddressByOrderCache(ctx, response)

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Retrieve active shipping addresses
// @Tags Shipping Address Query
// @Description Retrieve a list of active shipping addresses
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddressDeleteAt "List of active shipping addresses"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address-query/active [get]
func (h *shippingAddressQueryHandleApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetShippingAddressActiveCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByActive")
	}

	response := h.mapper.ToApiResponsePaginationShippingAddressDeleteAt(res)
	h.cache.SetShippingAddressActiveCache(ctx, req, response)

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Retrieve trashed shipping addresses
// @Tags Shipping Address Query
// @Description Retrieve a list of trashed shipping address records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddressDeleteAt "List of trashed shipping address data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address-query/trashed [get]
func (h *shippingAddressQueryHandleApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	search := c.QueryParam("search")

	req := &requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	ctx := c.Request().Context()
	if cached, found := h.cache.GetShippingAddressTrashedCache(ctx, req); found {
		return c.JSON(http.StatusOK, cached)
	}

	grpcReq := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "FindByTrashed")
	}

	response := h.mapper.ToApiResponsePaginationShippingAddressDeleteAt(res)
	h.cache.SetShippingAddressTrashedCache(ctx, req, response)

	return c.JSON(http.StatusOK, response)
}

func (h *shippingAddressQueryHandleApi) handleGrpcError(err error, operation string) error {
	h.logger.Error("Failed to "+operation, zap.Error(err))
	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to "+operation)
}
