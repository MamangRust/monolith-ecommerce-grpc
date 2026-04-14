package transactionhandler

import (
	"net/http"
	"strconv"

	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/cache/transaction"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/transaction"
	"github.com/labstack/echo/v4"
)

type transactionStatsHandlerApi struct {
	statsClient           pb.TransactionStatsServiceClient
	statsByMerchantClient pb.TransactionStatsByMerchantServiceClient
	logger                logger.LoggerInterface
	statsMapper           apimapper.TransactionStatsResponseMapper
	statsCache            transaction_cache.TransactionStatsCache
	statsByMerchantCache  transaction_cache.TransactionStatsByMerchantCache
}

type transactionStatsHandleDeps struct {
	statsClient           pb.TransactionStatsServiceClient
	statsByMerchantClient pb.TransactionStatsByMerchantServiceClient
	router                *echo.Echo
	logger                logger.LoggerInterface
	statsMapper           apimapper.TransactionStatsResponseMapper
	statsCache            transaction_cache.TransactionStatsCache
	statsByMerchantCache  transaction_cache.TransactionStatsByMerchantCache
}

func NewTransactionStatsHandleApi(params *transactionStatsHandleDeps) *transactionStatsHandlerApi {
	handler := &transactionStatsHandlerApi{
		statsClient:           params.statsClient,
		statsByMerchantClient: params.statsByMerchantClient,
		logger:                params.logger,
		statsMapper:           params.statsMapper,
		statsCache:            params.statsCache,
		statsByMerchantCache:  params.statsByMerchantCache,
	}

	routerTransaction := params.router.Group("/api/transaction-stats")

	// Stats
	routerTransaction.GET("/monthly-success", handler.FindMonthStatusSuccess)
	routerTransaction.GET("/yearly-success", handler.FindYearStatusSuccess)
	routerTransaction.GET("/monthly-failed", handler.FindMonthStatusFailed)
	routerTransaction.GET("/yearly-failed", handler.FindYearStatusFailed)

	routerTransaction.GET("/merchant/monthly-success", handler.FindMonthStatusSuccessByMerchant)
	routerTransaction.GET("/merchant/yearly-success", handler.FindYearStatusSuccessByMerchant)
	routerTransaction.GET("/merchant/monthly-failed", handler.FindMonthStatusFailedByMerchant)
	routerTransaction.GET("/merchant/yearly-failed", handler.FindYearStatusFailedByMerchant)

	routerTransaction.GET("/monthly-method-success", handler.FindMonthMethodSuccess)
	routerTransaction.GET("/yearly-method-success", handler.FindYearMethodSuccess)
	routerTransaction.GET("/merchant/monthly-method-success", handler.FindMonthMethodByMerchantSuccess)
	routerTransaction.GET("/merchant/yearly-method-success", handler.FindYearMethodByMerchantSuccess)

	routerTransaction.GET("/monthly-method-failed", handler.FindMonthMethodFailed)
	routerTransaction.GET("/yearly-method-failed", handler.FindYearMethodFailed)
	routerTransaction.GET("/merchant/monthly-method-failed", handler.FindMonthMethodByMerchantFailed)
	routerTransaction.GET("/merchant/yearly-method-failed", handler.FindYearMethodByMerchantFailed)

	return handler
}

// @Security Bearer
// @Summary Get monthly successful transaction amount
// @Tags Transaction Stats
// @Description Retrieve monthly successful transaction amount and count
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess "Monthly successful transaction amount"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/monthly-success [get]
func (h *transactionStatsHandlerApi) FindMonthStatusSuccess(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	ctx := c.Request().Context()
	req := &requests.MonthAmountTransaction{Year: year, Month: month}

	if cachedData, found := h.statsCache.GetCachedMonthAmountSuccessCached(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetMonthlyAmountSuccess(ctx, &pb.MonthAmountTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthAmountSuccess(res)
	h.statsCache.SetCachedMonthAmountSuccessCached(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly successful transaction amount
// @Tags Transaction Stats
// @Description Retrieve yearly successful transaction amount and count
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess "Yearly successful transaction amount"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/yearly-success [get]
func (h *transactionStatsHandlerApi) FindYearStatusSuccess(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	ctx := c.Request().Context()

	if cachedData, found := h.statsCache.GetCachedYearAmountSuccessCached(ctx, year); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetYearlyAmountSuccess(ctx, &pb.YearAmountTransactionRequest{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearAmountSuccess(res)
	h.statsCache.SetCachedYearAmountSuccessCached(ctx, year, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly failed transaction amount
// @Tags Transaction Stats
// @Description Retrieve monthly failed transaction amount and count
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed "Monthly failed transaction amount"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/monthly-failed [get]
func (h *transactionStatsHandlerApi) FindMonthStatusFailed(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	ctx := c.Request().Context()
	req := &requests.MonthAmountTransaction{Year: year, Month: month}

	if cachedData, found := h.statsCache.GetCachedMonthAmountFailedCached(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetMonthlyAmountFailed(ctx, &pb.MonthAmountTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthAmountFailed(res)
	h.statsCache.SetCachedMonthAmountFailedCached(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly failed transaction amount
// @Tags Transaction Stats
// @Description Retrieve yearly failed transaction amount and count
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed "Yearly failed transaction amount"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/yearly-failed [get]
func (h *transactionStatsHandlerApi) FindYearStatusFailed(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	ctx := c.Request().Context()

	if cachedData, found := h.statsCache.GetCachedYearAmountFailedCached(ctx, year); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetYearlyAmountFailed(ctx, &pb.YearAmountTransactionRequest{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearAmountFailed(res)
	h.statsCache.SetCachedYearAmountFailedCached(ctx, year, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly successful transaction amount by merchant
// @Tags Transaction Stats
// @Description Retrieve monthly successful transaction amount for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess "Monthly successful transaction amount by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/monthly-success [get]
func (h *transactionStatsHandlerApi) FindMonthStatusSuccessByMerchant(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	merchantID, _ := strconv.Atoi(c.QueryParam("merchant_id"))
	ctx := c.Request().Context()
	req := &requests.MonthAmountTransactionMerchant{Year: year, Month: month, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedMonthAmountSuccessByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetMonthlyAmountSuccessByMerchant(ctx, &pb.MonthAmountTransactionMerchantRequest{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthAmountSuccess(res)
	h.statsByMerchantCache.SetCachedMonthAmountSuccessByMerchant(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly successful transaction amount by merchant
// @Tags Transaction Stats
// @Description Retrieve yearly successful transaction amount for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess "Yearly successful transaction amount by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/yearly-success [get]
func (h *transactionStatsHandlerApi) FindYearStatusSuccessByMerchant(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	merchantID, _ := strconv.Atoi(c.QueryParam("merchant_id"))
	ctx := c.Request().Context()
	req := &requests.YearAmountTransactionMerchant{Year: year, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedYearAmountSuccessByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetYearlyAmountSuccessByMerchant(ctx, &pb.YearAmountTransactionMerchantRequest{
		Year: int32(year), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearAmountSuccess(res)
	h.statsByMerchantCache.SetCachedYearAmountSuccessByMerchant(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly failed transaction amount by merchant
// @Tags Transaction Stats
// @Description Retrieve monthly failed transaction amount for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed "Monthly failed transaction amount by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/monthly-failed [get]
func (h *transactionStatsHandlerApi) FindMonthStatusFailedByMerchant(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	merchantID, _ := strconv.Atoi(c.QueryParam("merchant_id"))
	ctx := c.Request().Context()
	req := &requests.MonthAmountTransactionMerchant{Year: year, Month: month, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedMonthAmountFailedByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetMonthlyAmountFailedByMerchant(ctx, &pb.MonthAmountTransactionMerchantRequest{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthAmountFailed(res)
	h.statsByMerchantCache.SetCachedMonthAmountFailedByMerchant(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly failed transaction amount by merchant
// @Tags Transaction Stats
// @Description Retrieve yearly failed transaction amount for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed "Yearly failed transaction amount by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/yearly-failed [get]
func (h *transactionStatsHandlerApi) FindYearStatusFailedByMerchant(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	merchantID, _ := strconv.Atoi(c.QueryParam("merchant_id"))
	ctx := c.Request().Context()
	req := &requests.YearAmountTransactionMerchant{Year: year, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedYearAmountFailedByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetYearlyAmountFailedByMerchant(ctx, &pb.YearAmountTransactionMerchantRequest{
		Year: int32(year), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearAmountFailed(res)
	h.statsByMerchantCache.SetCachedYearAmountFailedByMerchant(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly successful transaction method stats
// @Tags Transaction Stats
// @Description Retrieve monthly successful transaction method statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod "Monthly successful transaction method stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/monthly-method-success [get]
func (h *transactionStatsHandlerApi) FindMonthMethodSuccess(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	ctx := c.Request().Context()
	req := &requests.MonthMethodTransaction{Year: year, Month: month}

	if cachedData, found := h.statsCache.GetCachedMonthMethodSuccessCached(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetMonthlyTransactionMethodSuccess(ctx, &pb.MonthMethodTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthMethod(res)
	h.statsCache.SetCachedMonthMethodSuccessCached(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly successful transaction method stats
// @Tags Transaction Stats
// @Description Retrieve yearly successful transaction method statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod "Yearly successful transaction method stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/yearly-method-success [get]
func (h *transactionStatsHandlerApi) FindYearMethodSuccess(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	ctx := c.Request().Context()

	if cachedData, found := h.statsCache.GetCachedYearMethodSuccessCached(ctx, year); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetYearlyTransactionMethodSuccess(ctx, &pb.YearMethodTransactionRequest{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearMethod(res)
	h.statsCache.SetCachedYearMethodSuccessCached(ctx, year, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly successful transaction method stats by merchant
// @Tags Transaction Stats
// @Description Retrieve monthly successful transaction method statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod "Monthly successful transaction method stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/monthly-method-success [get]
func (h *transactionStatsHandlerApi) FindMonthMethodByMerchantSuccess(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	merchantID, _ := strconv.Atoi(c.QueryParam("merchant_id"))
	ctx := c.Request().Context()
	req := &requests.MonthMethodTransactionMerchant{Year: year, Month: month, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedMonthMethodSuccessByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetMonthlyTransactionMethodByMerchantSuccess(ctx, &pb.MonthMethodTransactionMerchantRequest{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthMethod(res)
	h.statsByMerchantCache.SetCachedMonthMethodSuccessByMerchant(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly successful transaction method stats by merchant
// @Tags Transaction Stats
// @Description Retrieve yearly successful transaction method statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod "Yearly successful transaction method stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/yearly-method-success [get]
func (h *transactionStatsHandlerApi) FindYearMethodByMerchantSuccess(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	merchantID, _ := strconv.Atoi(c.QueryParam("merchant_id"))
	ctx := c.Request().Context()
	req := &requests.YearMethodTransactionMerchant{Year: year, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedYearMethodSuccessByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetYearlyTransactionMethodByMerchantSuccess(ctx, &pb.YearMethodTransactionMerchantRequest{
		Year: int32(year), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearMethod(res)
	h.statsByMerchantCache.SetCachedYearMethodSuccessByMerchant(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly failed transaction method stats
// @Tags Transaction Stats
// @Description Retrieve monthly failed transaction method statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod "Monthly failed transaction method stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/monthly-method-failed [get]
func (h *transactionStatsHandlerApi) FindMonthMethodFailed(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	ctx := c.Request().Context()
	req := &requests.MonthMethodTransaction{Year: year, Month: month}

	if cachedData, found := h.statsCache.GetCachedMonthMethodFailedCached(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetMonthlyTransactionMethodFailed(ctx, &pb.MonthMethodTransactionRequest{Year: int32(year), Month: int32(month)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthMethod(res)
	h.statsCache.SetCachedMonthMethodFailedCached(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly failed transaction method stats
// @Tags Transaction Stats
// @Description Retrieve yearly failed transaction method statistics
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod "Yearly failed transaction method stats"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/yearly-method-failed [get]
func (h *transactionStatsHandlerApi) FindYearMethodFailed(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	ctx := c.Request().Context()

	if cachedData, found := h.statsCache.GetCachedYearMethodFailedCached(ctx, year); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsClient.GetYearlyTransactionMethodFailed(ctx, &pb.YearMethodTransactionRequest{Year: int32(year)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearMethod(res)
	h.statsCache.SetCachedYearMethodFailedCached(ctx, year, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get monthly failed transaction method stats by merchant
// @Tags Transaction Stats
// @Description Retrieve monthly failed transaction method statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param month query int false "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod "Monthly failed transaction method stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/monthly-method-failed [get]
func (h *transactionStatsHandlerApi) FindMonthMethodByMerchantFailed(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	month, _ := strconv.Atoi(c.QueryParam("month"))
	merchantID, _ := strconv.Atoi(c.QueryParam("merchant_id"))
	ctx := c.Request().Context()
	req := &requests.MonthMethodTransactionMerchant{Year: year, Month: month, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedMonthMethodFailedByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetMonthlyTransactionMethodByMerchantFailed(ctx, &pb.MonthMethodTransactionMerchantRequest{
		Year: int32(year), Month: int32(month), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionMonthMethod(res)
	h.statsByMerchantCache.SetCachedMonthMethodFailedByMerchant(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Get yearly failed transaction method stats by merchant
// @Tags Transaction Stats
// @Description Retrieve yearly failed transaction method statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int false "Year"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod "Yearly failed transaction method stats by merchant"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/transaction-stats/merchant/yearly-method-failed [get]
func (h *transactionStatsHandlerApi) FindYearMethodByMerchantFailed(c echo.Context) error {
	year, _ := strconv.Atoi(c.QueryParam("year"))
	merchantID, _ := strconv.Atoi(c.QueryParam("merchant_id"))
	ctx := c.Request().Context()
	req := &requests.YearMethodTransactionMerchant{Year: year, MerchantID: merchantID}

	if cachedData, found := h.statsByMerchantCache.GetCachedYearMethodFailedByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.statsByMerchantClient.GetYearlyTransactionMethodByMerchantFailed(ctx, &pb.YearMethodTransactionMerchantRequest{
		Year: int32(year), MerchantId: int32(merchantID),
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.statsMapper.ToApiResponseTransactionYearMethod(res)
	h.statsByMerchantCache.SetCachedYearMethodFailedByMerchant(ctx, req, apiResponse)
	return c.JSON(http.StatusOK, apiResponse)
}

