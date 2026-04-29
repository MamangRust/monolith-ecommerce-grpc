package shippingaddresshandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	shippingaddress_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/shipping_address"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/shipping_address"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingAddressCommandHandleApi struct {
	client pb.ShippingCommandServiceClient
	logger logger.LoggerInterface
	mapper apimapper.ShippingAddressCommandResponseMapper
	cache  shippingaddress_cache.ShippingAddressCommandCache
}

type shippingAddressCommandHandleDeps struct {
	client pb.ShippingCommandServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.ShippingAddressCommandResponseMapper
	cache  shippingaddress_cache.ShippingAddressCommandCache
}

func NewShippingAddressCommandHandleApi(deps *shippingAddressCommandHandleDeps) {
	handler := &shippingAddressCommandHandleApi{
		client: deps.client,
		logger: deps.logger,
		mapper: deps.mapper,
		cache:  deps.cache,
	}

	router := deps.router.Group("/api/shipping-address-command")
	router.POST("/trashed/:id", handler.TrashedShippingAddress)
	router.POST("/restore/:id", handler.RestoreShippingAddress)
	router.DELETE("/permanent/:id", handler.DeleteShippingAddressPermanent)
	router.POST("/restore/all", handler.RestoreAllShippingAddress)
	router.POST("/permanent/all", handler.DeleteAllShippingAddressPermanent)
}

// @Security Bearer
// @Summary Move shipping address to trash
// @Tags Shipping Address Command
// @Description Move a shipping address record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDeleteAt "Successfully moved shipping address to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move shipping address to trash"
// @Router /api/shipping-address-command/trashed/{id} [post]
func (h *shippingAddressCommandHandleApi) TrashedShippingAddress(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	res, err := h.client.TrashedShipping(ctx, &pb.FindByIdShippingRequest{Id: int32(id)})
	if err != nil {
		return h.handleGrpcError(err, "Trash")
	}

	h.cache.DeleteShippingAddressCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseShippingAddressDeleteAt(res))
}

// @Security Bearer
// @Summary Restore a trashed shipping address
// @Tags Shipping Address Command
// @Description Restore a trashed shipping address record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDeleteAt "Successfully restored shipping address"
// @Failure 400 {object} errors.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore shipping address"
// @Router /api/shipping-address-command/restore/{id} [post]
func (h *shippingAddressCommandHandleApi) RestoreShippingAddress(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	res, err := h.client.RestoreShipping(ctx, &pb.FindByIdShippingRequest{Id: int32(id)})
	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	h.cache.DeleteShippingAddressCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseShippingAddressDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete a shipping address
// @Tags Shipping Address Command
// @Description Permanently delete a shipping address record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Shipping Address ID"
// @Success 200 {object} response.ApiResponseShippingAddressDelete "Successfully deleted shipping address record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid shipping address ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete shipping address permanently"
// @Router /api/shipping-address-command/permanent/{id} [delete]
func (h *shippingAddressCommandHandleApi) DeleteShippingAddressPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	res, err := h.client.DeleteShippingPermanent(ctx, &pb.FindByIdShippingRequest{Id: int32(id)})
	if err != nil {
		return h.handleGrpcError(err, "Delete")
	}

	h.cache.DeleteShippingAddressCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseShippingAddressDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed shipping addresses
// @Tags Shipping Address Command
// @Description Restore all trashed shipping address records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseShippingAddressAll "Successfully restored all shipping addresses"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore shipping addresses"
// @Router /api/shipping-address-command/restore/all [post]
func (h *shippingAddressCommandHandleApi) RestoreAllShippingAddress(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllShipping(ctx, &emptypb.Empty{})
	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseShippingAddressAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed shipping addresses
// @Tags Shipping Address Command
// @Description Permanently delete all trashed shipping address records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseShippingAddressAll "Successfully deleted all shipping addresses permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete shipping addresses permanently"
// @Router /api/shipping-address-command/permanent/all [post]
func (h *shippingAddressCommandHandleApi) DeleteAllShippingAddressPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllShippingPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseShippingAddressAll(res))
}

func (h *shippingAddressCommandHandleApi) handleGrpcError(err error, operation string) error {
	h.logger.Error("Failed to "+operation, zap.Error(err))
	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to "+operation)
}
