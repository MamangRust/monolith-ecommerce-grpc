package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingAddressHandleApi struct {
	client  pb.ShippingServiceClient
	logger  logger.LoggerInterface
	mapping response_api.ShippingAddressResponseMapper
}

func NewHandlerShippingAddress(
	router *echo.Echo,
	client pb.ShippingServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ShippingAddressResponseMapper,
) *shippingAddressHandleApi {
	shippingHandler := &shippingAddressHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	myrouter := router.Group("/api/shipping-address")

	myrouter.GET("", shippingHandler.FindAllShipping)
	myrouter.GET("/:id", shippingHandler.FindById)
	myrouter.GET("/order/:id", shippingHandler.FindByOrder)
	myrouter.GET("/active", shippingHandler.FindByActive)
	myrouter.GET("/trashed", shippingHandler.FindByTrashed)

	myrouter.POST("/trashed/:id", shippingHandler.TrashedShippingAddress)
	myrouter.POST("/restore/:id", shippingHandler.RestoreShippingAddress)
	router.DELETE("/permanent/:id", shippingHandler.DeleteShippingAddressPermanent)

	myrouter.POST("/restore/all", shippingHandler.RestoreAllShippingAddress)
	myrouter.POST("/permanent/all", shippingHandler.DeleteAllShippingAddressPermanent)

	return shippingHandler

}

// @Security Bearer
// @Summary Find all shipping-address
// @Tags shipping address
// @Description Retrieve a list of all shipping-address
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddress "List of shipping-address"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping-address data"
// @Router /api/shipping-address [get]
func (h *shippingAddressHandleApi) FindAllShipping(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch shipping address", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedFindAllShippingAddresses(c)
	}

	so := h.mapping.ToApiResponsePaginationShippingAddress(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find shipping address by ID
// @Tags ShippingAddress
// @Description Retrieve a shipping address by ID
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddress "Shipping address data"
// @Failure 400 {object} response.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address/{id} [get]
func (h *shippingAddressHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid shipping address ID", zap.Error(err))
		return shippingaddress_errors.ErrApiInvalidIdShippingAddress(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch shipping address details", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedFindShippingAddressById(c)
	}

	so := h.mapping.ToApiResponseShippingAddress(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find shipping address by order ID
// @Tags ShippingAddress
// @Description Retrieve a shipping address by order ID
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.ApiResponseShippingAddress "Shipping address data"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping address data"
// @Router /api/shipping-address/order/{order_id} [get]
func (h *shippingAddressHandleApi) FindByOrder(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("order_id"))

	if err != nil {
		h.logger.Debug("Invalid order ID", zap.Error(err))
		return shippingaddress_errors.ErrApiInvalidOrderIdShippingAddress(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdShippingRequest{
		Id: int32(orderID),
	}

	res, err := h.client.FindByOrder(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch shipping address details", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedFindShippingAddressByOrder(c)
	}

	so := h.mapping.ToApiResponseShippingAddress(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active shipping-address
// @Tags ShippingAddress
// @Description Retrieve a list of active shipping-address
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddressDeleteAt "List of active shipping-address"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping data"
// @Router /api/shipping-address/active [get]
func (h *shippingAddressHandleApi) FindByActive(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch active shipping address", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedFindActiveShippingAddresses(c)
	}

	so := h.mapping.ToApiResponsePaginationShippingAddressDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed shipping-address records.
// @Summary Retrieve trashed shipping-address
// @Tags ShippingAddress
// @Description Retrieve a list of trashed shipping-address records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationShippingAddressDeleteAt "List of trashed shipping-address data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve shipping-address data"
// @Router /api/shipping-address/trashed [get]
func (h *shippingAddressHandleApi) FindByTrashed(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllShippingRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch archived shipping address", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedFindTrashedShippingAddresses(c)
	}

	so := h.mapping.ToApiResponsePaginationShippingAddressDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedShippingAddress retrieves a trashed shipping address record by its ID.
// @Summary Retrieve a trashed shipping address
// @Tags ShippingAddress
// @Description Retrieve a trashed shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDeleteAt "Successfully retrieved trashed shipping address"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed shipping address"
// @Router /api/shipping-address/trashed/{id} [get]
func (h *shippingAddressHandleApi) TrashedShippingAddress(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid shipping address ID format", zap.Error(err))
		return shippingaddress_errors.ErrApiInvalidIdShippingAddress(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedShipping(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive shipping address", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedTrashShippingAddress(c)
	}

	so := h.mapping.ToApiResponseShippingAddressDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreShippingAddress restores a shipping address record from the trash by its ID.
// @Summary Restore a trashed shipping address
// @Tags ShippingAddress
// @Description Restore a trashed shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDeleteAt "Successfully restored shipping address"
// @Failure 400 {object} response.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore shipping address"
// @Router /api/shipping-address/restore/{id} [post]
func (h *shippingAddressHandleApi) RestoreShippingAddress(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid shipping address ID", zap.Error(err))
		return shippingaddress_errors.ErrApiInvalidIdShippingAddress(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreShipping(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to restore shipping address", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedRestoreShippingAddress(c)
	}

	so := h.mapping.ToApiResponseShippingAddressDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteShippingAddressPermanent permanently deletes a shipping address record by its ID.
// @Summary Permanently delete a shipping address
// @Tags ShippingAddress
// @Description Permanently delete a shipping address record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDelete "Successfully deleted shipping address record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete shipping address:"
// @Router /api/shipping-address/delete/{id} [delete]
func (h *shippingAddressHandleApi) DeleteShippingAddressPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid shipping address ID", zap.Error(err))
		return shippingaddress_errors.ErrApiInvalidIdShippingAddress(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdShippingRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteShippingPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to permanently delete shipping address", zap.Error(err))

		return shippingaddress_errors.ErrApiFailedDeleteShippingAddressPermanent(c)
	}

	so := h.mapping.ToApiResponseShippingAddressDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllShippingAddress restores all trashed shipping address records.
// @Summary Restore all trashed shipping addresses
// @Tags ShippingAddress
// @Description Restore all trashed shipping address records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseShippingAddressAll "Successfully restored all shipping addresses"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all shipping addresses"
// @Router /api/shipping-address/restore/all [post]
func (h *shippingAddressHandleApi) RestoreAllShippingAddress(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllShipping(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk shipping address restoration failed", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedRestoreAllShippingAddresses(c)
	}

	so := h.mapping.ToApiResponseShippingAddressAll(res)

	h.logger.Debug("Successfully restored all shipping addresses")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllShippingAddressPermanent permanently deletes all trashed shipping address records.
// @Summary Permanently delete all trashed shipping addresses
// @Tags ShippingAddress
// @Description Permanently delete all trashed shipping address records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseShippingAddressAll "Successfully deleted all shipping addresses permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete all shipping addresses"
// @Router /api/shipping-address/delete/all [post]
func (h *shippingAddressHandleApi) DeleteAllShippingAddressPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllShippingPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk shipping address deletion failed", zap.Error(err))
		return shippingaddress_errors.ErrApiFailedDeleteAllPermanentShippingAddresses(c)
	}

	so := h.mapping.ToApiResponseShippingAddressAll(res)

	h.logger.Debug("Successfully deleted all shipping addresses permanently")

	return c.JSON(http.StatusOK, so)
}
