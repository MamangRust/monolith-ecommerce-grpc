package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/service"
// Removed unused db import
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type transactionStatsByMerchantHandler struct {
	pb.UnimplementedTransactionStatsByMerchantServiceServer
	service service.TransactionStatsByMerchantService
	logger  logger.LoggerInterface
}

func NewTransactionStatsByMerchantHandler(service service.TransactionStatsByMerchantService, logger logger.LoggerInterface) *transactionStatsByMerchantHandler {
	return &transactionStatsByMerchantHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionStatsByMerchantHandler) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *pb.MonthAmountTransactionMerchantRequest) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	request := &requests.MonthAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountSuccessByMerchant(ctx, request)
	if err != nil {
		return nil, err
	}

	var stats []*pb.TransactionMonthlyAmountSuccess
	for _, v := range data {
		stats = append(stats, &pb.TransactionMonthlyAmountSuccess{
			Year:         v.Year,
			Month:        v.Month,
			TotalSuccess: int32(v.TotalSuccess),
			TotalAmount:  int32(v.TotalAmount),
		})
	}

	return &pb.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched monthly amount success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *pb.YearAmountTransactionMerchantRequest) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	request := &requests.YearAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyAmountSuccessByMerchant(ctx, request)
	if err != nil {
		return nil, err
	}

	var stats []*pb.TransactionYearlyAmountSuccess
	for _, v := range data {
		stats = append(stats, &pb.TransactionYearlyAmountSuccess{
			Year:         v.Year,
			TotalSuccess: int32(v.TotalSuccess),
			TotalAmount:  int32(v.TotalAmount),
		})
	}

	return &pb.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Successfully fetched yearly amount success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *pb.MonthAmountTransactionMerchantRequest) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	request := &requests.MonthAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountFailedByMerchant(ctx, request)
	if err != nil {
		return nil, err
	}

	var stats []*pb.TransactionMonthlyAmountFailed
	for _, v := range data {
		stats = append(stats, &pb.TransactionMonthlyAmountFailed{
			Year:        v.Year,
			Month:       v.Month,
			TotalFailed: int32(v.TotalFailed),
			TotalAmount: int32(v.TotalAmount),
		})
	}

	return &pb.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Successfully fetched monthly amount failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyAmountFailedByMerchant(ctx context.Context, req *pb.YearAmountTransactionMerchantRequest) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	request := &requests.YearAmountTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyAmountFailedByMerchant(ctx, request)
	if err != nil {
		return nil, err
	}

	var stats []*pb.TransactionYearlyAmountFailed
	for _, v := range data {
		stats = append(stats, &pb.TransactionYearlyAmountFailed{
			Year:        v.Year,
			TotalFailed: int32(v.TotalFailed),
			TotalAmount: int32(v.TotalAmount),
		})
	}

	return &pb.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Successfully fetched yearly amount failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *pb.MonthMethodTransactionMerchantRequest) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodByMerchantSuccess(ctx, request)
	if err != nil {
		return nil, err
	}

	var stats []*pb.TransactionMonthlyMethod
	for _, v := range data {
		stats = append(stats, &pb.TransactionMonthlyMethod{
			Month:             v.Month,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *pb.YearMethodTransactionMerchantRequest) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	request := &requests.YearMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyMethodByMerchantSuccess(ctx, request)
	if err != nil {
		return nil, err
	}

	var stats []*pb.TransactionYearlyMethod
	for _, v := range data {
		stats = append(stats, &pb.TransactionYearlyMethod{
			Year:              v.Year,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method success stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *pb.MonthMethodTransactionMerchantRequest) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
		Month:      int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodByMerchantFailed(ctx, request)
	if err != nil {
		return nil, err
	}

	var stats []*pb.TransactionMonthlyMethod
	for _, v := range data {
		stats = append(stats, &pb.TransactionMonthlyMethod{
			Month:             v.Month,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Successfully fetched monthly transaction method failed stats by merchant",
		Data:    stats,
	}, nil
}

func (h *transactionStatsByMerchantHandler) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *pb.YearMethodTransactionMerchantRequest) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	request := &requests.YearMethodTransactionMerchant{
		MerchantID: int(req.GetMerchantId()),
		Year:       int(req.GetYear()),
	}

	data, err := h.service.FindYearlyMethodByMerchantFailed(ctx, request)
	if err != nil {
		return nil, err
	}

	var stats []*pb.TransactionYearlyMethod
	for _, v := range data {
		stats = append(stats, &pb.TransactionYearlyMethod{
			Year:              v.Year,
			PaymentMethod:     v.PaymentMethod,
			TotalTransactions: int32(v.TotalTransactions),
			TotalAmount:       int32(v.TotalAmount),
		})
	}

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Successfully fetched yearly transaction method failed stats by merchant",
		Data:    stats,
	}, nil
}
