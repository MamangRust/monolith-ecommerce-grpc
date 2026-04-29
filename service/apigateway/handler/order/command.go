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
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderCommandHandlerApi struct {
	client pb.OrderCommandServiceClient
	logger logger.LoggerInterface
	mapper apimapper.OrderCommandResponseMapper
	cache  order_cache.OrderCommandCache
}

type orderCommandHandleDeps struct {
	client pb.OrderCommandServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.OrderCommandResponseMapper
	cache  order_cache.OrderCommandCache
}

func NewOrderCommandHandleApi(params *orderCommandHandleDeps) *orderCommandHandlerApi {
	handler := &orderCommandHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
		cache:  params.cache,
	}

	routerOrder := params.router.Group("/api/order-command")
	routerOrder.POST("/create", handler.Create)
	routerOrder.POST("/update/:id", handler.Update)
	routerOrder.POST("/trashed/:id", handler.Trash)
	routerOrder.POST("/restore/:id", handler.Restore)
	routerOrder.DELETE("/permanent/:id", handler.DeletePermanent)
	routerOrder.POST("/restore/all", handler.RestoreAll)
	routerOrder.POST("/permanent/all", handler.DeleteAllPermanent)

	return handler
}

// @Security Bearer
// @Summary Create a new order
// @Tags Order Command
// @Description Create a new order with items and shipping address
// @Accept json
// @Produce json
// @Param request body requests.CreateOrderRequest true "Order details"
// @Success 200 {object} response.ApiResponseOrder "Successfully created order"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 500 {object} errors.ErrorResponse "Failed to create order"
// @Router /api/order-command/create [post]
func (h *orderCommandHandlerApi) Create(c echo.Context) error {
	var body requests.CreateOrderRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	
	items := make([]*pb.CreateOrderItemRequest, 0)
	for _, item := range body.Items {
		items = append(items, &pb.CreateOrderItemRequest{
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     int32(item.Price),
		})
	}

	res, err := h.client.Create(ctx, &pb.CreateOrderRequest{
		MerchantId: int32(body.MerchantID),
		UserId:     int32(body.UserID),
		TotalPrice: int32(body.TotalPrice),
		Items:      items,
		Shipping: &pb.CreateShippingAddressRequest{
			Alamat:         body.ShippingAddress.Alamat,
			Provinsi:       body.ShippingAddress.Provinsi,
			Kota:           body.ShippingAddress.Kota,
			Courier:        body.ShippingAddress.Courier,
			ShippingMethod: body.ShippingAddress.ShippingMethod,
			ShippingCost:   int32(body.ShippingAddress.ShippingCost),
			Negara:         body.ShippingAddress.Negara,
		},
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseOrder(res))
}

// @Security Bearer
// @Summary Update an existing order
// @Tags Order Command
// @Description Update an existing order record
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body requests.UpdateOrderRequest true "Updated order details"
// @Success 200 {object} response.ApiResponseOrder "Successfully updated order"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 500 {object} errors.ErrorResponse "Failed to update order"
// @Router /api/order-command/update/{id} [post]
func (h *orderCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	var body requests.UpdateOrderRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()

	items := make([]*pb.UpdateOrderItemRequest, 0)
	for _, item := range body.Items {
		items = append(items, &pb.UpdateOrderItemRequest{
			OrderItemId: int32(item.OrderItemID),
			ProductId:   int32(item.ProductID),
			Quantity:    int32(item.Quantity),
			Price:       int32(item.Price),
		})
	}

	res, err := h.client.Update(ctx, &pb.UpdateOrderRequest{
		OrderId:    int32(id),
		UserId:     int32(body.UserID),
		TotalPrice: int32(body.TotalPrice),
		Items:      items,
		Shipping: &pb.UpdateShippingAddressRequest{
			ShippingId:     0, // Will be set by service
			Alamat:         body.ShippingAddress.Alamat,
			Provinsi:       body.ShippingAddress.Provinsi,
			Kota:           body.ShippingAddress.Kota,
			Courier:        body.ShippingAddress.Courier,
			ShippingMethod: body.ShippingAddress.ShippingMethod,
			ShippingCost:   int32(body.ShippingAddress.ShippingCost),
			Negara:         body.ShippingAddress.Negara,
		},
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteOrderCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseOrder(res))
}

// @Security Bearer
// @Summary Move order to trash
// @Tags Order Command
// @Description Move an order record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDeleteAt "Successfully moved order to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid order ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move order to trash"
// @Router /api/order-command/trashed/{id} [post]
func (h *orderCommandHandlerApi) Trash(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.TrashedOrder(ctx, &pb.FindByIdOrderRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteOrderCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseOrderDeleteAt(res))
}

// @Security Bearer
// @Summary Restore a trashed order
// @Tags Order Command
// @Description Restore a trashed order record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDeleteAt "Successfully restored order"
// @Failure 400 {object} errors.ErrorResponse "Invalid order ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore order"
// @Router /api/order-command/restore/{id} [post]
func (h *orderCommandHandlerApi) Restore(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.RestoreOrder(ctx, &pb.FindByIdOrderRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteOrderCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseOrderDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete an order
// @Tags Order Command
// @Description Permanently delete an order record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} response.ApiResponseOrderDelete "Successfully deleted order record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid order ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete order permanently"
// @Router /api/order-command/permanent/{id} [delete]
func (h *orderCommandHandlerApi) DeletePermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.DeleteOrderPermanent(ctx, &pb.FindByIdOrderRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteOrderCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseOrderDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed orders
// @Tags Order Command
// @Description Restore all trashed order records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseOrderAll "Successfully restored all orders"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore orders"
// @Router /api/order-command/restore/all [post]
func (h *orderCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllOrder(ctx, &emptypb.Empty{})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteOrderCache(ctx, 0)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseOrderAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed orders
// @Tags Order Command
// @Description Permanently delete all trashed order records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseOrderAll "Successfully deleted all orders permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete orders permanently"
// @Router /api/order-command/permanent/all [post]
func (h *orderCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllOrderPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteOrderCache(ctx, 0)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseOrderAll(res))
}

