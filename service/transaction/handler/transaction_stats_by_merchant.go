package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-transaction/service"
	// Removed unused db import
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"

	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
)

type transactionStatsByMerchantHandler struct {
	pbtransaction.UnimplementedTransactionStatsByMerchantServiceServer
	service service.TransactionStatsByMerchantService
	logger  logger.LoggerInterface
}

func NewTransactionStatsByMerchantHandler(service service.TransactionStatsByMerchantService, logger logger.LoggerInterface) *transactionStatsByMerchantHandler {
	return &transactionStatsByMerchantHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionStatsByMerchantHandler) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *pbtransaction.MonthAmountTransactionMerchantRequest) (*pbtransaction.ApiResponseTransactionMonthAmountSuccess, error) {
	request := &requests.MonthAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountSuccessByMerchant(ctx, request)
	if err != nil {
		return nil, err
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
		Message: "Successfully fetched monthly amount success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *pbtransaction.YearAmountTransactionMerchantRequest) (*pbtransaction.ApiResponseTransactionYearAmountSuccess, error) {
	request := &requests.YearAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyAmountSuccessByMerchant(ctx, request)
	if err != nil {
		return nil, err
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
		Message: "Successfully fetched yearly amount success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *pbtransaction.MonthAmountTransactionMerchantRequest) (*pbtransaction.ApiResponseTransactionMonthAmountFailed, error) {
	request := &requests.MonthAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountFailedByMerchant(ctx, request)
	if err != nil {
		return nil, err
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
		Message: "Successfully fetched monthly amount failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyAmountFailedByMerchant(ctx context.Context, req *pbtransaction.YearAmountTransactionMerchantRequest) (*pbtransaction.ApiResponseTransactionYearAmountFailed, error) {
	request := &requests.YearAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyAmountFailedByMerchant(ctx, request)
	if err != nil {
		return nil, err
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
		Message: "Successfully fetched yearly amount failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *pbtransaction.MonthMethodTransactionMerchantRequest) (*pbtransaction.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodByMerchantSuccess(ctx, request)
	if err != nil {
		return nil, err
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
		Message: "Successfully fetched monthly transaction method success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *pbtransaction.YearMethodTransactionMerchantRequest) (*pbtransaction.ApiResponseTransactionYearPaymentmethod, error) {
	request := &requests.YearMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyMethodByMerchantSuccess(ctx, request)
	if err != nil {
		return nil, err
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
		Message: "Successfully fetched yearly transaction method success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *pbtransaction.MonthMethodTransactionMerchantRequest) (*pbtransaction.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodByMerchantFailed(ctx, request)
	if err != nil {
		return nil, err
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
		Message: "Successfully fetched monthly transaction method failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *pbtransaction.YearMethodTransactionMerchantRequest) (*pbtransaction.ApiResponseTransactionYearPaymentmethod, error) {
	request := &requests.YearMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyMethodByMerchantFailed(ctx, request)
	if err != nil {
		return nil, err
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
		Message: "Successfully fetched yearly transaction method failed stats by merchant",
		Data:    stats,
	}, nil
}
