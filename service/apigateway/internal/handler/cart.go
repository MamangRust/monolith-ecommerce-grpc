package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/cart_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type cartHandleApi struct {
	client  pb.CartServiceClient
	logger  logger.LoggerInterface
	mapping response_api.CartResponseMapper
}

func NewHandlerCart(
	router *echo.Echo,
	client pb.CartServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CartResponseMapper,
) *cartHandleApi {
	cartHandler := &cartHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routerCart := router.Group("/api/cart")
	routerCart.GET("", cartHandler.FindAll)
	routerCart.POST("/create", cartHandler.Create)
	routerCart.DELETE("/:id", cartHandler.Delete)
	routerCart.POST("/delete-all", cartHandler.DeleteAll)

	return cartHandler
}

// @Security Bearer
// @Summary Find all carts
// @Tags Cart
// @Description Retrieve a list of all carts
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponseCartPagination "List of carts"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve cart data"
// @Router /api/cart [get]
func (h *cartHandleApi) FindAll(c echo.Context) error {
	userID, err := strconv.Atoi(c.QueryParam("user_id"))

	if err != nil || userID <= 0 {
		h.logger.Debug("Invalid user ID format", zap.Error(err))
		return cart_errors.ErrApiFailedInvalidUserId(c)
	}

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

	req := &pb.FindAllCartRequest{
		UserId:   int32(userID),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch cart details", zap.Error(err))
		return cart_errors.ErrApiFailedFindAllCarts(c)
	}

	so := h.mapping.ToApiResponseCartPagination(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new cart
// @Tags Cart
// @Description Create a new cart item
// @Accept json
// @Produce json
// @Param body body requests.CreateCartRequest true "Cart creation data"
// @Success 200 {object} response.ApiResponseCart "Created cart details"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 502 {object} response.ErrorResponse "Failed to create cart"
// @Router /api/cart [post]
func (h *cartHandleApi) Create(c echo.Context) error {
	var body requests.CreateCartRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return cart_errors.ErrApiBindCreateCart(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return cart_errors.ErrApiValidateCreateCart(c)
	}

	ctx := c.Request().Context()
	req := &pb.CreateCartRequest{
		Quantity:  int32(body.Quantity),
		ProductId: int32(body.ProductID),
		UserId:    int32(body.UserID),
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		h.logger.Error("cart creation failed", zap.Error(err))
		return cart_errors.ErrApiFailedCreateCart(c)
	}

	so := h.mapping.ToApiResponseCart(res)
	return c.JSON(http.StatusCreated, so)
}

// @Security Bearer
// @Summary Delete a cart
// @Tags Cart
// @Description Delete a cart by ID
// @Accept json
// @Produce json
// @Param id path int true "Cart ID"
// @Success 200 {object} response.ApiResponseCartDelete "Successfully deleted cart"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete cart"
// @Router /api/cart/{id} [delete]
func (h *cartHandleApi) Delete(c echo.Context) error {
	id := c.Param("id")

	idStr, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))

		return cart_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCartRequest{
		Id: int32(idStr),
	}

	res, err := h.client.Delete(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete cart", zap.Error(err))
		return cart_errors.ErrApiFailedDeleteCart(c)
	}

	so := h.mapping.ToApiResponseCartDelete(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Delete multiple carts
// @Tags Cart
// @Description Delete multiple carts by IDs
// @Accept json
// @Produce json
// @Param request body requests.DeleteCartRequest true "Cart IDs"
// @Success 200 {object} response.ApiResponseCartAll "Successfully deleted carts"
// @Failure 500 {object} response.ErrorResponse "Failed to delete carts"
// @Router /api/cart/delete-all [post]
func (h *cartHandleApi) DeleteAll(c echo.Context) error {
	var req requests.DeleteCartRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Debug("Invalid request parameter", zap.Error(err))
		return cart_errors.ErrApiBindDeleteAllCart(c)
	}

	if err := req.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return cart_errors.ErrApiValidateDeleteAllCart(c)
	}

	ctx := c.Request().Context()

	cartIdsPb := make([]int32, len(req.CartIds))
	for i, id := range req.CartIds {
		cartIdsPb[i] = int32(id)
	}

	reqPb := &pb.DeleteCartRequest{
		CartIds: cartIdsPb,
	}

	res, err := h.client.DeleteAll(ctx, reqPb)

	if err != nil {
		h.logger.Error("Failed to delete cart items", zap.Error(err))
		return cart_errors.ErrApiFailedDeleteAllCarts(c)
	}

	so := h.mapping.ToApiResponseCartAll(res)
	return c.JSON(http.StatusOK, so)
}
