package carthandler

import (
	"net/http"

	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/cart"
	"github.com/labstack/echo/v4"
)

type cartCommandHandlerApi struct {
	client pb.CartCommandServiceClient
	logger logger.LoggerInterface
	mapper apimapper.CartCommandResponseMapper
}

type cartCommandHandleDeps struct {
	client pb.CartCommandServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.CartCommandResponseMapper
}

func NewCartCommandHandleApi(params *cartCommandHandleDeps) *cartCommandHandlerApi {
	handler := &cartCommandHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
	}

	routerCart := params.router.Group("/api/cart-query")
	routerCart.POST("/create", handler.Create)
	routerCart.DELETE("/delete", handler.Delete)
	routerCart.POST("/delete-all", handler.DeleteAll)

	return handler
}

// @Security Bearer
// @Summary Add item to cart
// @Tags Cart Command
// @Description Add a product to the user's cart
// @Accept json
// @Produce json
// @Param request body requests.CreateCartRequest true "Cart item details"
// @Success 201 {object} response.ApiResponseCart "Successfully added to cart"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 500 {object} errors.ErrorResponse "Failed to add to cart"
// @Router /api/cart-query/create [post]
func (h *cartCommandHandlerApi) Create(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok || userID <= 0 { return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized") }

	var body requests.CreateCartRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, &pb.CreateCartRequest{
		UserId:    int32(userID),
		ProductId: int32(body.ProductID),
		Quantity:  int32(body.Quantity),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusCreated, h.mapper.ToApiResponseCart(res))
}

// @Security Bearer
// @Summary Remove item from cart
// @Tags Cart Command
// @Description Remove a specific item from the user's cart
// @Accept json
// @Produce json
// @Param request body requests.DeleteCartRequest true "Cart ID to delete"
// @Success 200 {object} response.ApiResponseCartDelete "Successfully removed from cart"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Failed to remove from cart"
// @Router /api/cart-query/delete [delete]
func (h *cartCommandHandlerApi) Delete(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok || userID <= 0 { return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized") }

	var body requests.DeleteCartRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	res, err := h.client.Delete(ctx, &pb.DeleteCartRequest{
		UserId: int32(userID),
		CartId: int32(body.CartID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseCartDelete(res))
}

// @Security Bearer
// @Summary Remove multiple items from cart
// @Tags Cart Command
// @Description Remove multiple specific items from the user's cart
// @Accept json
// @Produce json
// @Param request body requests.DeleteAllCartRequest true "Cart IDs to delete"
// @Success 200 {object} response.ApiResponseCartAll "Successfully removed all items from cart"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Failed to remove items from cart"
// @Router /api/cart-query/delete-all [post]
func (h *cartCommandHandlerApi) DeleteAll(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok || userID <= 0 { return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized") }

	var body requests.DeleteAllCartRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	cartIdsPb := make([]int32, len(body.CartIds))
	for i, id := range body.CartIds { cartIdsPb[i] = int32(id) }

	res, err := h.client.DeleteAll(ctx, &pb.DeleteAllCartRequest{
		UserId:  int32(userID),
		CartIds: cartIdsPb,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseCartAll(res))
}

