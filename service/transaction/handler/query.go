package handler

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-transaction/service"

	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
)

type transactionQueryHandler struct {
	pbtransaction.UnimplementedTransactionQueryServiceServer
	service service.TransactionQueryService
	logger  logger.LoggerInterface
}

func NewTransactionQueryHandler(service service.TransactionQueryService, logger logger.LoggerInterface) *transactionQueryHandler {
	return &transactionQueryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionQueryHandler) FindAllTransactions(ctx context.Context, req *pbtransaction.FindAllTransactionRequest) (*pbtransaction.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindAll(ctx, request)
	if err != nil {
		return nil, err
	}

	var transactions []*pbtransaction.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponse(v))
	}

	return &pbtransaction.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindByActive(ctx context.Context, req *pbtransaction.FindAllTransactionRequest) (*pbtransaction.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindActive(ctx, request)
	if err != nil {
		return nil, err
	}

	var transactions []*pbtransaction.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseActive(v))
	}

	return &pbtransaction.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched active transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindByTrashed(ctx context.Context, req *pbtransaction.FindAllTransactionRequest) (*pbtransaction.ApiResponsePaginationTransactionDeleteAt, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindTrashed(ctx, request)
	if err != nil {
		return nil, err
	}

	var transactions []*pbtransaction.TransactionResponseDeleteAt
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseDeleteAt(v))
	}

	return &pbtransaction.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindById(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	data, err := h.service.FindByID(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction",
		Data:    h.ToTransactionResponseId(data),
	}, nil
}

func (h *transactionQueryHandler) FindByOrderId(ctx context.Context, req *pbtransaction.FindByOrderIdTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	data, err := h.service.FindByOrderID(ctx, int(req.GetOrderId()))
	if err != nil {
		return nil, err
	}

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction by order id",
		Data:    h.ToTransactionResponseOrderId(data),
	}, nil
}

func (h *transactionQueryHandler) FindByMerchant(ctx context.Context, req *pbtransaction.FindAllTransactionByMerchantRequest) (*pbtransaction.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransactionByMerchant{
		MerchantID: int(req.GetMerchantId()),
		Page:       int(req.GetPage()),
		PageSize:   int(req.GetPageSize()),
		Search:     req.GetSearch(),
	}

	data, total, err := h.service.FindByMerchant(ctx, request)
	if err != nil {
		return nil, err
	}

	var transactions []*pbtransaction.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseMerchant(v))
	}

	return &pbtransaction.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transactions by merchant",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

// Manual Mappings

func (h *transactionQueryHandler) ToTransactionResponse(v *db.GetTransactionsRow) *pbtransaction.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseActive(v *db.GetTransactionsActiveRow) *pbtransaction.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseDeleteAt(v *db.GetTransactionsTrashedRow) *pbtransaction.TransactionResponseDeleteAt {
	return mapToProtoTransactionResponseDeleteAt(v)
}

func (h *transactionQueryHandler) ToTransactionResponseId(v *db.GetTransactionByIDRow) *pbtransaction.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseOrderId(v *db.GetTransactionByOrderIDRow) *pbtransaction.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseMerchant(v *db.GetTransactionByMerchantRow) *pbtransaction.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}
