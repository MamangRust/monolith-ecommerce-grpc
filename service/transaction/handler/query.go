package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type transactionQueryHandler struct {
	pb.UnimplementedTransactionQueryServiceServer
	service service.TransactionQueryService
	logger  logger.LoggerInterface
}

func NewTransactionQueryHandler(service service.TransactionQueryService, logger logger.LoggerInterface) *transactionQueryHandler {
	return &transactionQueryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionQueryHandler) FindAllTransactions(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindAll(ctx, request)
	if err != nil {
		return nil, err
	}

	var transactions []*pb.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponse(v))
	}

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindByActive(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindActive(ctx, request)
	if err != nil {
		return nil, err
	}

	var transactions []*pb.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseActive(v))
	}

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched active transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindByTrashed(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	request := &requests.FindAllTransaction{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	}

	data, total, err := h.service.FindTrashed(ctx, request)
	if err != nil {
		return nil, err
	}

	var transactions []*pb.TransactionResponseDeleteAt
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseDeleteAt(v))
	}

	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed transactions",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

func (h *transactionQueryHandler) FindById(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	data, err := h.service.FindByID(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction",
		Data:    h.ToTransactionResponseId(data),
	}, nil
}

func (h *transactionQueryHandler) FindByOrderId(ctx context.Context, req *pb.FindByOrderIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	data, err := h.service.FindByOrderID(ctx, int(req.GetOrderId()))
	if err != nil {
		return nil, err
	}

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction by order id",
		Data:    h.ToTransactionResponseOrderId(data),
	}, nil
}

func (h *transactionQueryHandler) FindByMerchant(ctx context.Context, req *pb.FindAllTransactionByMerchantRequest) (*pb.ApiResponsePaginationTransaction, error) {
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

	var transactions []*pb.TransactionResponse
	for _, v := range data {
		transactions = append(transactions, h.ToTransactionResponseMerchant(v))
	}

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transactions by merchant",
		Data:       transactions,
		Pagination: createPaginationMeta(request.Page, request.PageSize, *total),
	}, nil
}

// Manual Mappings

func (h *transactionQueryHandler) ToTransactionResponse(v *db.GetTransactionsRow) *pb.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseActive(v *db.GetTransactionsActiveRow) *pb.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseDeleteAt(v *db.GetTransactionsTrashedRow) *pb.TransactionResponseDeleteAt {
	return mapToProtoTransactionResponseDeleteAt(v)
}

func (h *transactionQueryHandler) ToTransactionResponseId(v *db.GetTransactionByIDRow) *pb.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseOrderId(v *db.GetTransactionByOrderIDRow) *pb.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionQueryHandler) ToTransactionResponseMerchant(v *db.GetTransactionByMerchantRow) *pb.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}
