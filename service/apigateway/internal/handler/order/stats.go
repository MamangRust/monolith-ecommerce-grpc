package orderhandler

import (
	"net/http"
	"strconv"

	order_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/cache/order"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/order"
	"github.com/labstack/echo/v4"
)

type orderStatsHandlerApi struct {
	client            pb.OrderQueryServiceClient
	logger            logger.LoggerInterface
	mapper            apimapper.OrderStatsResponseMapper
	cache             order_cache.OrderStatsCache
	merchantStatsCache order_cache.OrderStatsByMerchantCache
}

type orderStatsHandleDeps struct {
	client            pb.OrderQueryServiceClient
	router            *echo.Echo
	logger            logger.LoggerInterface
	mapper            apimapper.OrderStatsResponseMapper
	cache             order_cache.OrderStatsCache
	merchantStatsCache order_cache.OrderStatsByMerchantCache
}

func NewOrderStatsHandleApi(params *orderStatsHandleDeps) *orderStatsHandlerApi {
	handler := &orderStatsHandlerApi{
		client:            params.client,
		logger:            params.logger,
		mapper:            params.mapper,
		cache:             params.cache,
		merchantStatsCache: params.merchantStatsCache,
	}

	routerOrder := params.router.Group("/api/order")
	routerOrder.GET("/monthly-total-revenue", handler.FindMonthlyTotalRevenue)
	routerOrder.GET("/yearly-total-revenue", handler.FindYearlyTotalRevenue)
	routerOrder.GET("/merchant/monthly-total-revenue", handler.FindMonthlyTotalRevenueByMerchant)
	routerOrder.GET("/merchant/yearly-total-revenue", handler.FindYearlyTotalRevenueByMerchant)

	routerOrder.GET("/monthly-revenue", handler.FindMonthlyRevenue)
	routerOrder.GET("/yearly-revenue", handler.FindYearlyRevenue)
	routerOrder.GET("/merchant/monthly-revenue", handler.FindMonthlyRevenueByMerchant)
	routerOrder.GET("/merchant/yearly-revenue", handler.FindYearlyRevenueByMerchant)

	return handler
}

// @Security Bearer
// @Summary Get monthly total revenue
// @Tags Order Stats
// @Description Retrieve monthly total revenue and order stats
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponseOrderMonthlyTotalRevenue "Monthly total revenue"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/order/monthly-total-revenue [get]
func (h *orderStatsHandlerApi) FindMonthlyTotalRevenue(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))

	ctx := c.Request().Context()
	req := &requests.MonthTotalRevenue{Year: year, Month: month}

	if cachedData, found := h.cache.GetMonthlyTotalRevenueCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTotalRevenue(ctx, &pb.FindYearMonthTotalRevenue{
		Year: int32(year), Month: int32(month),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseMonthlyTotalRevenue(res)
	h.cache.SetMonthlyTotalRevenueCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly total revenue
// @Tags Order Stats
// @Description Retrieve yearly total revenue and order stats
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponseOrderYearlyTotalRevenue "Yearly total revenue"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/order/yearly-total-revenue [get]
func (h *orderStatsHandlerApi) FindYearlyTotalRevenue(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetYearlyTotalRevenueCache(ctx, year); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTotalRevenue(ctx, &pb.FindYearTotalRevenue{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseYearlyTotalRevenue(res)
	h.cache.SetYearlyTotalRevenueCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly total revenue by merchant
// @Tags Order Stats
// @Description Retrieve monthly total revenue for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderMonthlyTotalRevenue "Monthly total revenue by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/order/merchant/monthly-total-revenue [get]
func (h *orderStatsHandlerApi) FindMonthlyTotalRevenueByMerchant(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	merchantId, _ := strconv.Atoi(c.QueryParam("merchant_id"))

	ctx := c.Request().Context()
	req := &requests.MonthTotalRevenueMerchant{Year: year, Month: month, MerchantID: merchantId}

	if cachedData, found := h.merchantStatsCache.GetMonthlyTotalRevenueByMerchantCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyTotalRevenueByMerchant(ctx, &pb.FindYearMonthTotalRevenueByMerchant{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantId),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseMonthlyTotalRevenue(res)
	h.merchantStatsCache.SetMonthlyTotalRevenueByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly total revenue by merchant
// @Tags Order Stats
// @Description Retrieve yearly total revenue for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderYearlyTotalRevenue "Yearly total revenue by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/order/merchant/yearly-total-revenue [get]
func (h *orderStatsHandlerApi) FindYearlyTotalRevenueByMerchant(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	merchantId, _ := strconv.Atoi(c.QueryParam("merchant_id"))

	ctx := c.Request().Context()
	req := &requests.YearTotalRevenueMerchant{Year: year, MerchantID: merchantId}

	if cachedData, found := h.merchantStatsCache.GetYearlyTotalRevenueByMerchantCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyTotalRevenueByMerchant(ctx, &pb.FindYearTotalRevenueByMerchant{
		Year: int32(year), MerchantId: int32(merchantId),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseYearlyTotalRevenue(res)
	h.merchantStatsCache.SetYearlyTotalRevenueByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly revenue stats
// @Tags Order Stats
// @Description Retrieve monthly revenue statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponseOrderMonthly "Monthly revenue stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/order/monthly-revenue [get]
func (h *orderStatsHandlerApi) FindMonthlyRevenue(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetMonthlyOrderCache(ctx, year); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyRevenue(ctx, &pb.FindYearOrder{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseMonthlyOrder(res)
	h.cache.SetMonthlyOrderCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly revenue stats
// @Tags Order Stats
// @Description Retrieve yearly revenue statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponseOrderYearly "Yearly revenue stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/order/yearly-revenue [get]
func (h *orderStatsHandlerApi) FindYearlyRevenue(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetYearlyOrderCache(ctx, year); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyRevenue(ctx, &pb.FindYearOrder{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseYearlyOrder(res)
	h.cache.SetYearlyOrderCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly revenue stats by merchant
// @Tags Order Stats
// @Description Retrieve monthly revenue statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderMonthly "Monthly revenue stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/order/merchant/monthly-revenue [get]
func (h *orderStatsHandlerApi) FindMonthlyRevenueByMerchant(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	merchantId, _ := strconv.Atoi(c.QueryParam("merchant_id"))

	ctx := c.Request().Context()
	req := &requests.MonthOrderMerchant{Year: year, MerchantID: merchantId}

	if cachedData, found := h.merchantStatsCache.GetMonthlyOrderByMerchantCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindMonthlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year: int32(year), MerchantId: int32(merchantId),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseMonthlyOrder(res)
	h.merchantStatsCache.SetMonthlyOrderByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly revenue stats by merchant
// @Tags Order Stats
// @Description Retrieve yearly revenue statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseOrderYearly "Yearly revenue stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/order/merchant/yearly-revenue [get]
func (h *orderStatsHandlerApi) FindYearlyRevenueByMerchant(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	merchantId, _ := strconv.Atoi(c.QueryParam("merchant_id"))

	ctx := c.Request().Context()
	req := &requests.YearOrderMerchant{Year: year, MerchantID: merchantId}

	if cachedData, found := h.merchantStatsCache.GetYearlyOrderByMerchantCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindYearlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year: int32(year), MerchantId: int32(merchantId),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseYearlyOrder(res)
	h.merchantStatsCache.SetYearlyOrderByMerchantCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

