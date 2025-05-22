package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionHandleApi struct {
	client  pb.TransactionServiceClient
	kafka   kafka.Kafka
	logger  logger.LoggerInterface
	mapping response_api.TransactionResponseMapper
}

func NewHandlerTransaction(
	router *echo.Echo,
	client pb.TransactionServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.TransactionResponseMapper,
) *transactionHandleApi {
	transactionHandle := &transactionHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routercategory := router.Group("/api/transaction")

	routercategory.GET("", transactionHandle.FindAllTransaction)

	routercategory.GET("/active", transactionHandle.FindByActive)
	routercategory.GET("/:id", transactionHandle.FindById)
	routercategory.GET("/merchant/:merchant_id", transactionHandle.FindByMerchant)
	routercategory.GET("/trashed", transactionHandle.FindByTrashed)

	routercategory.GET("/monthly-success", transactionHandle.FindMonthStatusSuccess)
	routercategory.GET("/yearly-success", transactionHandle.FindYearStatusSuccess)
	routercategory.GET("/monthly-failed", transactionHandle.FindMonthStatusFailed)
	routercategory.GET("/yearly-failed", transactionHandle.FindYearStatusFailed)

	routercategory.GET("/merchant/monthly-success", transactionHandle.FindMonthStatusSuccessByMerchant)
	routercategory.GET("/merchant/yearly-success", transactionHandle.FindYearStatusSuccessByMerchant)
	routercategory.GET("/merchant/monthly-failed", transactionHandle.FindMonthStatusFailedByMerchant)
	routercategory.GET("/merchant/yearly-failed", transactionHandle.FindYearStatusFailedByMerchant)

	routercategory.GET("/merchant/monthly-success", transactionHandle.FindMonthStatusSuccessByMerchant)
	routercategory.GET("/merchant/yearly-success", transactionHandle.FindYearStatusSuccessByMerchant)
	routercategory.GET("/merchant/monthly-failed", transactionHandle.FindMonthStatusFailedByMerchant)
	routercategory.GET("/merchant/yearly-failed", transactionHandle.FindYearStatusFailedByMerchant)

	routercategory.GET("/monthly-method-success", transactionHandle.FindMonthMethodSuccess)
	routercategory.GET("/yearly-method-success", transactionHandle.FindYearMethodSuccess)

	routercategory.GET("/merchant/monthly-method-success", transactionHandle.FindMonthMethodByMerchantSuccess)
	routercategory.GET("/merchant/yearly-method-success", transactionHandle.FindYearMethodByMerchantSuccess)

	routercategory.GET("/monthly-method-failed", transactionHandle.FindMonthMethodFailed)
	routercategory.GET("/yearly-method-failed", transactionHandle.FindYearMethodFailed)

	routercategory.GET("/merchant/monthly-method-failed", transactionHandle.FindMonthMethodByMerchantFailed)
	routercategory.GET("/merchant/yearly-method-failed", transactionHandle.FindYearMethodByMerchantFailed)

	routercategory.POST("/create", transactionHandle.Create)
	routercategory.POST("/update/:id", transactionHandle.Update)

	routercategory.POST("/trashed/:id", transactionHandle.TrashedTransaction)
	routercategory.POST("/restore/:id", transactionHandle.RestoreTransaction)
	routercategory.DELETE("/permanent/:id", transactionHandle.DeleteTransactionPermanent)

	routercategory.POST("/restore/all", transactionHandle.RestoreAllTransaction)
	routercategory.POST("/permanent/all", transactionHandle.DeleteAllTransactionPermanent)

	return transactionHandle
}

// @Security Bearer
// @Summary Find all transactions
// @Tags Transaction
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction [get]
func (h *transactionHandleApi) FindAllTransaction(c echo.Context) error {
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

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch transactions", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationTransaction(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find all transactions by merchant
// @Tags Transaction
// @Description Retrieve a list of all transactions filtered by merchant
// @Accept json
// @Produce json
// @Param merchant_id path int true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/merchant/{merchant_id} [get]
func (h *transactionHandleApi) FindByMerchant(c echo.Context) error {
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))

	if err != nil || merchantID <= 0 {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
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

	req := &pb.FindAllTransactionMerchantRequest{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch transaction", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindByMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationTransaction(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find transaction by ID
// @Tags Transaction
// @Description Retrieve a transaction by ID
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/{id} [get]
func (h *transactionHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid transaction ID format", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch transaction details", zap.Error(err))
		return transaction_errors.ErrApiTransactionNotFound(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active transactions
// @Tags Transaction
// @Description Retrieve a list of active transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransactionDeleteAt "List of active transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/active [get]
func (h *transactionHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch active transactions", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed transaction records.
// @Summary Retrieve trashed transactions
// @Tags Transaction
// @Description Retrieve a list of trashed transaction records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransactionDeleteAt "List of trashed transaction data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/trashed [get]
func (h *transactionHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch archived transactions", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthStatusSuccess retrieves monthly successful transactions
// @Summary Get monthly successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-success [get]
func (h *transactionHandleApi) FindMonthStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthStatusSuccess(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction status success", zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthSuccess(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearStatusSuccess retrieves yearly successful transactions
// @Summary Get yearly successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-success [get]
func (h *transactionHandleApi) FindYearStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearStatusSuccess(ctx, &pb.FindYearlyTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction status success", zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearSuccess(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthStatusFailed retrieves monthly failed transactions
// @Summary Get monthly failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-failed [get]
func (h *transactionHandleApi) FindMonthStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthStatusFailed(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction status failed", zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthFailed(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearStatusFailed retrieves yearly failed transactions
// @Summary Get yearly failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-failed [get]
func (h *transactionHandleApi) FindYearStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearStatusFailed(ctx, &pb.FindYearlyTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction status failed", zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearFailed(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthStatusSuccessByMerchant retrieves monthly successful transactions by merchant
// @Summary Get monthly successful transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID, year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-success [get]
func (h *transactionHandleApi) FindMonthStatusSuccessByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthStatusSuccessByMerchant(ctx, &pb.FindMonthlyTransactionStatusByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction status success", zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthSuccessByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearStatusSuccessByMerchant retrieves yearly successful transactions by merchant
// @Summary Get yearly successful transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-success [get]
func (h *transactionHandleApi) FindYearStatusSuccessByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearStatusSuccessByMerchant(ctx, &pb.FindYearlyTransactionStatusByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction status success", zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearSuccessByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthStatusFailedByMerchant retrieves monthly failed transactions by merchant
// @Summary Get monthly failed transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID, year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-failed [get]
func (h *transactionHandleApi) FindMonthStatusFailedByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthStatusFailedByMerchant(ctx, &pb.FindMonthlyTransactionStatusByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction status failed", zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthFailedByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearStatusFailedByMerchant retrieves yearly failed transactions by merchant
// @Summary Get yearly failed transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-failed [get]
func (h *transactionHandleApi) FindYearStatusFailedByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearStatusFailedByMerchant(ctx, &pb.FindYearlyTransactionStatusByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction status failed", zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearFailedByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthMethod retrieves monthly payment method statistics
// @Summary Get monthly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-method-success [get]
func (h *transactionHandleApi) FindMonthMethodSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := strconv.Atoi(c.QueryParam("month"))

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthMethodSuccess(ctx, &pb.MonthTransactionMethod{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction methods", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindMonthMethod(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearMethod retrieves yearly payment method statistics
// @Summary Get yearly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-method-success [get]
func (h *transactionHandleApi) FindYearMethodSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearMethodSuccess(ctx, &pb.YearTransactionMethod{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction methods", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindYearMethod(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthMethodByMerchant retrieves monthly payment method statistics by merchant
// @Summary Get monthly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-method-success/{merchant_id} [get]
func (h *transactionHandleApi) FindMonthMethodByMerchantSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := strconv.Atoi(c.QueryParam("month"))

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthMethodByMerchantSuccess(ctx, &pb.MonthTransactionMethodByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly transaction methods", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindMonthMethodByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearMethodByMerchant retrieves yearly payment method statistics by merchant
// @Summary Get yearly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-method-success/{merchant_id} [get]
func (h *transactionHandleApi) FindYearMethodByMerchantSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearMethodByMerchantSuccess(ctx, &pb.YearTransactionMethodByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly transaction methods", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindYearMethodByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthMethod retrieves monthly payment method statistics
// @Summary Get monthly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-method-failed [get]
func (h *transactionHandleApi) FindMonthMethodFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := strconv.Atoi(c.QueryParam("month"))

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthMethodFailed(ctx, &pb.MonthTransactionMethod{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Success to retrieve monthly transaction methods", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindMonthMethod(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearMethod retrieves yearly payment method statistics
// @Summary Get yearly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-method-failed [get]
func (h *transactionHandleApi) FindYearMethodFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearMethodFailed(ctx, &pb.YearTransactionMethod{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Success to retrieve yearly transaction methods", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindYearMethod(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthMethodByMerchant retrieves monthly payment method statistics by merchant
// @Summary Get monthly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-method-failed/{merchant_id} [get]
func (h *transactionHandleApi) FindMonthMethodByMerchantFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := strconv.Atoi(c.QueryParam("month"))

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthMethodByMerchantSuccess(ctx, &pb.MonthTransactionMethodByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		h.logger.Debug("Success to retrieve monthly transaction methods", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindMonthMethodByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearMethodByMerchant retrieves yearly payment method statistics by merchant
// @Summary Get yearly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-method-failed/{merchant_id} [get]
func (h *transactionHandleApi) FindYearMethodByMerchantFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearMethodByMerchantSuccess(ctx, &pb.YearTransactionMethodByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Success to retrieve yearly transaction methods", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindYearMethodByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new transaction
// @Tags Transaction
// @Description Create a new transaction record
// @Accept json
// @Produce json
// @Param request body requests.CreateTransactionRequest true "Transaction details"
// @Success 200 {object} response.ApiResponseTransaction "Successfully created transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transaction"
// @Router /api/transaction/create [post]
func (h *transactionHandleApi) Create(c echo.Context) error {
	userID := c.Get("userID").(string)

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		h.logger.Debug("Invalid user ID", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidUserId(c)
	}

	var req requests.CreateTransactionRequest

	req.UserID = userIDInt

	if err := c.Bind(&req); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return transaction_errors.ErrApiBindCreateTransaction(c)
	}

	if err := req.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		log.Fatal(err)
		return transaction_errors.ErrApiValidateCreateTransaction(c)
	}

	ctx := c.Request().Context()
	grpcReq := &pb.CreateTransactionRequest{
		OrderId:       int32(req.OrderID),
		MerchantId:    int32(req.MerchantID),
		PaymentMethod: req.PaymentMethod,
		Amount:        int32(req.Amount),
	}

	res, err := h.client.Create(ctx, grpcReq)
	if err != nil {
		h.logger.Error("transaction creation failed", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedCreate(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update an existing transaction
// @Tags Transaction
// @Description Update an existing transaction record
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body requests.UpdateTransactionRequest true "Updated transaction details"
// @Success 200 {object} response.ApiResponseTransaction "Successfully updated transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transaction"
// @Router /api/transaction/update [post]
func (h *transactionHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	var req requests.UpdateTransactionRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))

		return transaction_errors.ErrApiBindUpdateTransaction(c)
	}

	if err := req.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return transaction_errors.ErrApiValidateUpdateTransaction(c)
	}

	ctx := c.Request().Context()
	grpcReq := &pb.UpdateTransactionRequest{
		TransactionId: int32(idInt),
		OrderId:       int32(req.OrderID),
		MerchantId:    int32(req.MerchantID),
		PaymentMethod: req.PaymentMethod,
		Amount:        int32(req.Amount),
	}

	res, err := h.client.Update(ctx, grpcReq)

	if err != nil {
		h.logger.Error("transaction update failed", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedTransaction retrieves a trashed transaction record by its ID.
// @Summary Retrieve a trashed transaction
// @Tags Transaction
// @Description Retrieve a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDeleteAt "Successfully retrieved trashed transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed transaction"
// @Router /api/transaction/trashed/{id} [get]
func (h *transactionHandleApi) TrashedTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid transaction ID format", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedTransaction(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive transaction", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedTrashed(c)
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreTransaction restores a transaction record from the trash by its ID.
// @Summary Restore a trashed transaction
// @Tags Transaction
// @Description Restore a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDeleteAt "Successfully restored transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction"
// @Router /api/transaction/restore/{id} [post]
func (h *transactionHandleApi) RestoreTransaction(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid transaction ID format", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreTransaction(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore transaction", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedRestore(c)
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteTransactionPermanent permanently deletes a transaction record by its ID.
// @Summary Permanently delete a transaction
// @Tags Transaction
// @Description Permanently delete a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDelete "Successfully deleted transaction record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction"
// @Router /api/transaction/delete/{id} [delete]
func (h *transactionHandleApi) DeleteTransactionPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid transaction ID format", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteTransactionPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete transaction", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseTransactionDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllTransaction restores all trashed transactions.
// @Summary Restore all trashed transactions
// @Tags Transaction
// @Description Restore all trashed transactions.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully restored all transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transactions"
// @Router /api/transaction/restore/all [post]
func (h *transactionHandleApi) RestoreAllTransaction(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllTransaction(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk transactions restoration failed", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	h.logger.Debug("Successfully restored all transactions")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllTransactionPermanent permanently deletes all transactions.
// @Summary Permanently delete all transactions
// @Tags Transaction
// @Description Permanently delete all transactions.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully deleted all transactions permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transactions"
// @Router /api/transaction/delete/all [post]
func (h *transactionHandleApi) DeleteAllTransactionPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk transactions deletion failed", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	h.logger.Debug("Successfully deleted all transactions permanently")

	return c.JSON(http.StatusOK, so)
}
