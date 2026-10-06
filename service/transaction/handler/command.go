package handler

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-transaction/service"
	"google.golang.org/protobuf/types/known/emptypb"

	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
)

type transactionCommandHandler struct {
	pbtransaction.UnimplementedTransactionCommandServiceServer
	service service.TransactionCommandService
	logger  logger.LoggerInterface
}

func NewTransactionCommandHandler(service service.TransactionCommandService, logger logger.LoggerInterface) *transactionCommandHandler {
	return &transactionCommandHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionCommandHandler) Create(ctx context.Context, req *pbtransaction.CreateTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	request := &requests.CreateTransactionRequest{
		OrderID:       int(req.GetOrderId()),
		MerchantID:    int(req.GetMerchantId()),
		UserID:        int(req.GetUserId()),
		PaymentMethod: req.GetPaymentMethod(),
		Amount:        int(req.GetAmount()),
		PaymentStatus: &[]string{req.GetPaymentStatus()}[0],
	}

	data, err := h.service.Create(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data:    h.ToTransactionResponseCreate(data),
	}, nil
}

func (h *transactionCommandHandler) Update(ctx context.Context, req *pbtransaction.UpdateTransactionRequest) (*pbtransaction.ApiResponseTransaction, error) {
	transactionID := int(req.GetTransactionId())
	request := &requests.UpdateTransactionRequest{
		TransactionID: &transactionID,
		MerchantID:    int(req.GetMerchantId()),
		OrderID:       int(req.GetOrderId()),
		PaymentMethod: req.GetPaymentMethod(),
		Amount:        int(req.GetAmount()),
	}

	if status := req.GetPaymentStatus(); status != "" {
		request.PaymentStatus = &status
	}

	data, err := h.service.Update(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data:    h.ToTransactionResponseUpdate(data),
	}, nil
}

func (h *transactionCommandHandler) TrashedTransaction(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDeleteAt, error) {
	data, err := h.service.Trash(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data:    h.ToTransactionResponseDeleteAt(data),
	}, nil
}

func (h *transactionCommandHandler) RestoreTransaction(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDeleteAt, error) {
	data, err := h.service.Restore(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data:    h.ToTransactionResponseDeleteAt(data),
	}, nil
}

func (h *transactionCommandHandler) DeleteTransactionPermanent(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDelete, error) {
	_, err := h.service.DeletePermanent(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transaction permanently",
	}, nil
}

func (h *transactionCommandHandler) RestoreAllTransaction(ctx context.Context, req *emptypb.Empty) (*pbtransaction.ApiResponseTransactionAll, error) {
	_, err := h.service.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restored all transactions",
	}, nil
}

func (h *transactionCommandHandler) DeleteTransactionByOrderPermanent(ctx context.Context, req *pbtransaction.FindByIdTransactionRequest) (*pbtransaction.ApiResponseTransactionDelete, error) {
	_, err := h.service.DeleteByOrderIDPermanent(ctx, int(req.GetId()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transactions by order permanently",
	}, nil
}

func (h *transactionCommandHandler) DeleteAllTransactionPermanent(ctx context.Context, req *emptypb.Empty) (*pbtransaction.ApiResponseTransactionAll, error) {
	_, err := h.service.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbtransaction.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully deleted all transactions permanently",
	}, nil
}

// Manual Mappings

func (h *transactionCommandHandler) ToTransactionResponseCreate(v *db.CreateTransactionRow) *pbtransaction.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionCommandHandler) ToTransactionResponseUpdate(v *db.UpdateTransactionRow) *pbtransaction.TransactionResponse {
	return mapToProtoTransactionResponse(v)
}

func (h *transactionCommandHandler) ToTransactionResponseDeleteAt(v *db.Transaction) *pbtransaction.TransactionResponseDeleteAt {
	return mapToProtoTransactionResponseDeleteAt(v)
}
