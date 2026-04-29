package transactionhandler

import (
	"net/http"
	"strconv"

	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/transaction"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/transaction"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionCommandHandlerApi struct {
	client     pb.TransactionCommandServiceClient
	logger     logger.LoggerInterface
	mapper     apimapper.TransactionCommandResponseMapper
	cache      transaction_cache.TransactionCommandCache
	apiHandler sharedErrors.ApiHandler
}

type transactionCommandHandleDeps struct {
	client     pb.TransactionCommandServiceClient
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     apimapper.TransactionCommandResponseMapper
	cache      transaction_cache.TransactionCommandCache
	apiHandler sharedErrors.ApiHandler
}

func NewTransactionCommandHandleApi(params *transactionCommandHandleDeps) *transactionCommandHandlerApi {
	handler := &transactionCommandHandlerApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	routerTransaction := params.router.Group("/api/transaction-command")
	routerTransaction.POST("/create", params.apiHandler.Handle("CreateTransaction", handler.Create))
	routerTransaction.POST("/update/:id", params.apiHandler.Handle("UpdateTransaction", handler.Update))
	routerTransaction.POST("/trashed/:id", params.apiHandler.Handle("TrashedTransaction", handler.Trashed))
	routerTransaction.POST("/restore/:id", params.apiHandler.Handle("RestoreTransaction", handler.Restore))
	routerTransaction.DELETE("/permanent/:id", params.apiHandler.Handle("DeleteTransactionPermanent", handler.DeletePermanent))
	routerTransaction.POST("/restore/all", params.apiHandler.Handle("RestoreAllTransaction", handler.RestoreAll))
	routerTransaction.POST("/permanent/all", params.apiHandler.Handle("DeleteAllTransactionPermanent", handler.DeleteAllPermanent))

	return handler
}

// @Security Bearer
// @Summary Create a new transaction
// @Tags Transaction Command
// @Description Create a new transaction with the provided details
// @Accept mpfd
// @Produce json
// @Param order_id formData int true "Order ID"
// @Param merchant_id formData int true "Merchant ID"
// @Param amount formData int true "Amount"
// @Param user_id formData int true "User ID"
// @Param status formData string true "Payment Status"
// @Param method formData string true "Payment Method"
// @Success 201 {object} response.ApiResponseTransaction "Successfully created transaction"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create transaction"
// @Router /api/transaction/create [post]
func (h *transactionCommandHandlerApi) Create(c echo.Context) error {
	req := new(pb.CreateTransactionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, req)
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.InvalidateTransactionCache(ctx)

	return c.JSON(http.StatusCreated, h.mapper.ToApiResponseTransaction(res))
}

// @Security Bearer
// @Summary Update transaction status
// @Tags Transaction Command
// @Description Update the payment status of a transaction
// @Accept mpfd
// @Produce json
// @Param id path int true "Transaction ID"
// @Param status formData string true "Payment Status"
// @Success 200 {object} response.ApiResponseTransaction "Successfully updated transaction"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update transaction"
// @Router /api/transaction/update/{id} [post]
func (h *transactionCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	status := c.FormValue("status")

	ctx := c.Request().Context()
	res, err := h.client.Update(ctx, &pb.UpdateTransactionRequest{
		TransactionId: int32(id),
		PaymentStatus: status,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteTransactionCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseTransaction(res))
}

// @Security Bearer
// @Summary Move transaction to trash
// @Tags Transaction Command
// @Description Move a transaction record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDeleteAt "Successfully moved transaction to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move transaction to trash"
// @Router /api/transaction/trashed/{id} [post]
func (h *transactionCommandHandlerApi) Trashed(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteTransactionCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseTransactionDeleteAt(res))
}

// @Security Bearer
// @Summary Restore a trashed transaction
// @Tags Transaction Command
// @Description Restore a trashed transaction record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDeleteAt "Successfully restored transaction"
// @Failure 400 {object} errors.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore transaction"
// @Router /api/transaction/restore/{id} [post]
func (h *transactionCommandHandlerApi) Restore(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.RestoreTransaction(ctx, &pb.FindByIdTransactionRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteTransactionCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseTransactionDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete a transaction
// @Tags Transaction Command
// @Description Permanently delete a transaction record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDelete "Successfully deleted transaction record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete transaction permanently"
// @Router /api/transaction/permanent/{id} [delete]
func (h *transactionCommandHandlerApi) DeletePermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.DeleteTransactionPermanent(ctx, &pb.FindByIdTransactionRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteTransactionCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseTransactionDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed transactions
// @Tags Transaction Command
// @Description Restore all trashed transaction records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully restored all transactions"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore transactions"
// @Router /api/transaction/restore/all [post]
func (h *transactionCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllTransaction(ctx, &emptypb.Empty{})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.InvalidateTransactionCache(ctx)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseTransactionAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed transactions
// @Tags Transaction Command
// @Description Permanently delete all trashed transaction records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully deleted all transactions permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete transactions permanently"
// @Router /api/transaction/permanent/all [post]
func (h *transactionCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.InvalidateTransactionCache(ctx)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseTransactionAll(res))
}

