package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-transaction/service"

	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
)

type transactionStatsHandler struct {
	pbtransaction.UnimplementedTransactionStatsServiceServer
	service service.TransactionStatsService
	logger  logger.LoggerInterface
}

func NewTransactionStatsHandler(service service.TransactionStatsService, logger logger.LoggerInterface) *transactionStatsHandler {
	return &transactionStatsHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionStatsHandler) GetMonthlyAmountSuccess(ctx context.Context, req *pbtransaction.MonthAmountTransactionRequest) (*pbtransaction.ApiResponseTransactionMonthAmountSuccess, error) {
	request := &requests.MonthAmountTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountSuccess(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransaction.TransactionMonthlyAmountSuccess
	for _, v := range data {
		stats = append(stats, &pbtransaction.TransactionMonthlyAmountSuccess{
			Year:         v.Year,
			Month:        v.Month,
			TotalSuccess: int32(v.TotalSuccess),
			TotalAmount:  int32(v.TotalAmount),
		})
	}

	return &pbtransaction.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly amount success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyAmountSuccess(ctx context.Context, req *pbtransaction.YearAmountTransactionRequest) (*pbtransaction.ApiResponseTransactionYearAmountSuccess, error) {
	data, err := h.service.FindYearlyAmountSuccess(ctx, int(req.GetYear()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransaction.TransactionYearlyAmountSuccess
	for _, v := range data {
		stats = append(stats, &pbtransaction.TransactionYearlyAmountSuccess{
			Year:         v.Year,
			TotalSuccess: int32(v.TotalSuccess),
			TotalAmount:  int32(v.TotalAmount),
		})
	}

	return &pbtransaction.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly amount success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyAmountFailed(ctx context.Context, req *pbtransaction.MonthAmountTransactionRequest) (*pbtransaction.ApiResponseTransactionMonthAmountFailed, error) {
	request := &requests.MonthAmountTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountFailed(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransaction.TransactionMonthlyAmountFailed
	for _, v := range data {
		stats = append(stats, &pbtransaction.TransactionMonthlyAmountFailed{
			Year:        v.Year,
			Month:       v.Month,
			TotalFailed: int32(v.TotalFailed),
			TotalAmount: int32(v.TotalAmount),
		})
	}

	return &pbtransaction.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Successfully fetched monthly amount failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyAmountFailed(ctx context.Context, req *pbtransaction.YearAmountTransactionRequest) (*pbtransaction.ApiResponseTransactionYearAmountFailed, error) {
	data, err := h.service.FindYearlyAmountFailed(ctx, int(req.GetYear()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransaction.TransactionYearlyAmountFailed
	for _, v := range data {
		stats = append(stats, &pbtransaction.TransactionYearlyAmountFailed{
			Year:        v.Year,
			TotalFailed: int32(v.TotalFailed),
			TotalAmount: int32(v.TotalAmount),
		})
	}

	return &pbtransaction.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Successfully fetched yearly amount failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *pbtransaction.MonthMethodTransactionRequest) (*pbtransaction.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodSuccess(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransaction.TransactionMonthlyMethod
	for _, v := range data {
		stats = append(stats, &pbtransaction.TransactionMonthlyMethod{
			Month:             v.Month,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pbtransaction.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyTransactionMethodSuccess(ctx context.Context, req *pbtransaction.YearMethodTransactionRequest) (*pbtransaction.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.service.FindYearlyMethodSuccess(ctx, int(req.GetYear()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransaction.TransactionYearlyMethod
	for _, v := range data {
		stats = append(stats, &pbtransaction.TransactionYearlyMethod{
			Year:              v.Year,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pbtransaction.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyTransactionMethodFailed(ctx context.Context, req *pbtransaction.MonthMethodTransactionRequest) (*pbtransaction.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodFailed(ctx, request)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransaction.TransactionMonthlyMethod
	for _, v := range data {
		stats = append(stats, &pbtransaction.TransactionMonthlyMethod{
			Month:             v.Month,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pbtransaction.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyTransactionMethodFailed(ctx context.Context, req *pbtransaction.YearMethodTransactionRequest) (*pbtransaction.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.service.FindYearlyMethodFailed(ctx, int(req.GetYear()))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var stats []*pbtransaction.TransactionYearlyMethod
	for _, v := range data {
		stats = append(stats, &pbtransaction.TransactionYearlyMethod{
			Year:              v.Year,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pbtransaction.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method failed stats",
		Data:    stats,
	}, nil
}
