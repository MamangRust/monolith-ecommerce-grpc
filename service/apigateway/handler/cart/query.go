package carthandler

import (
	"net/http"
	"strconv"

	cart_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/cart"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/cart"
	"github.com/labstack/echo/v4"
)

type cartQueryHandlerApi struct {
	client pb.CartQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.CartQueryResponseMapper
	cache  cart_cache.CartQueryCache
}

type cartQueryHandleDeps struct {
	client pb.CartQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.CartQueryResponseMapper
	cache  cart_cache.CartQueryCache
}

func NewCartQueryHandleApi(params *cartQueryHandleDeps) *cartQueryHandlerApi {
	handler := &cartQueryHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
		cache:  params.cache,
	}

	routerCart := params.router.Group("/api/cart-command")
	routerCart.GET("", handler.FindAll)

	return handler
}

// @Security Bearer
// @Summary Find all cart items
// @Tags Cart Query
// @Description Retrieve all items in the current user's cart
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponseCartPagination "List of cart items"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve cart data"
// @Router /api/cart-command [get]
func (h *cartQueryHandlerApi) FindAll(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok || userID <= 0 { return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized") }

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllCarts{UserID: userID, Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedCarts(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllCartRequest{
		UserId:   int32(userID),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseCartPagination(res)
	h.cache.SetCachedCarts(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

