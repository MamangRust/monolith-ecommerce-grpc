package transactionhandler

import (
	"net/http"
	"strconv"

	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/transaction"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/transaction"
	"github.com/labstack/echo/v4"
)

type transactionQueryHandlerApi struct {
	queryClient pb.TransactionQueryServiceClient
	logger      logger.LoggerInterface
	mapper      apimapper.TransactionQueryResponseMapper
	cache       transaction_cache.TransactionQueryCache
	apiHandler  sharedErrors.ApiHandler
}

type transactionQueryHandleDeps struct {
	queryClient pb.TransactionQueryServiceClient
	router      *echo.Echo
	logger      logger.LoggerInterface
	mapper      apimapper.TransactionQueryResponseMapper
	cache       transaction_cache.TransactionQueryCache
	apiHandler  sharedErrors.ApiHandler
}

func NewTransactionQueryHandleApi(params *transactionQueryHandleDeps) *transactionQueryHandlerApi {
	handler := &transactionQueryHandlerApi{
		queryClient: params.queryClient,
		logger:      params.logger,
		mapper:      params.mapper,
		cache:       params.cache,
		apiHandler:  params.apiHandler,
	}

	routerTransaction := params.router.Group("/api/transaction-query")
	routerTransaction.GET("", params.apiHandler.Handle("FindAllTransactions", handler.FindAll))
	routerTransaction.GET("/:id", params.apiHandler.Handle("FindById", handler.FindById))
	routerTransaction.GET("/merchant/:merchant_id", params.apiHandler.Handle("FindByMerchant", handler.FindByMerchant))
	routerTransaction.GET("/active", params.apiHandler.Handle("FindByActive", handler.FindByActive))
	routerTransaction.GET("/trashed", params.apiHandler.Handle("FindByTrashed", handler.FindByTrashed))

	return handler
}

// @Security Bearer
// @Summary Find all transactions
// @Tags Transaction Query
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query [get]
func (h *transactionQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllTransaction{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedTransactionsCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.queryClient.FindAllTransactions(ctx, &pb.FindAllTransactionRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationTransaction(res)
	h.cache.SetCachedTransactionsCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find transaction by ID
// @Tags Transaction Query
// @Description Retrieve a transaction by ID
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Transaction data"
// @Failure 400 {object} errors.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/{id} [get]
func (h *transactionQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid Transaction ID") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedTransactionCache(ctx, id); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.queryClient.FindById(ctx, &pb.FindByIdTransactionRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponseTransaction(res)
	h.cache.SetCachedTransactionCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find transactions by merchant ID
// @Tags Transaction Query
// @Description Retrieve a list of transactions belonging to a specific merchant
// @Accept json
// @Produce json
// @Param merchant_id path int true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions by merchant"
// @Failure 400 {object} errors.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/merchant/{merchant_id} [get]
func (h *transactionQueryHandlerApi) FindByMerchant(c echo.Context) error {
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid Merchant ID") }

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllTransactionByMerchant{MerchantID: merchantID, Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedTransactionByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.queryClient.FindByMerchant(ctx, &pb.FindAllTransactionByMerchantRequest{
		MerchantId: int32(merchantID), Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationTransaction(res)
	h.cache.SetCachedTransactionByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active transactions
// @Tags Transaction Query
// @Description Retrieve a list of active transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of active transactions"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/active [get]
func (h *transactionQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllTransaction{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedTransactionActiveCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.queryClient.FindByActive(ctx, &pb.FindAllTransactionRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationTransaction(res)
	h.cache.SetCachedTransactionActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed transactions
// @Tags Transaction Query
// @Description Retrieve a list of trashed transaction records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransactionDeleteAt "List of trashed transaction data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/trashed [get]
func (h *transactionQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllTransaction{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedTransactionTrashedCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.queryClient.FindByTrashed(ctx, &pb.FindAllTransactionRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	apiResponse := h.mapper.ToApiResponsePaginationTransactionDeleteAt(res)
	h.cache.SetCachedTransactionTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}
