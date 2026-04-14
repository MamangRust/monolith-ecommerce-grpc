package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/service"
// Removed unused db import
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type transactionStatsHandler struct {
	pb.UnimplementedTransactionStatsServiceServer
	service service.TransactionStatsService
	logger  logger.LoggerInterface
}

func NewTransactionStatsHandler(service service.TransactionStatsService, logger logger.LoggerInterface) *transactionStatsHandler {
	return &transactionStatsHandler{
		service: service,
		logger:  logger,
	}
}

func (h *transactionStatsHandler) GetMonthlyAmountSuccess(ctx context.Context, req *pb.MonthAmountTransactionRequest) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	request := &requests.MonthAmountTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountSuccess(ctx, request)
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
		Message: "Successfully fetched monthly amount success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyAmountSuccess(ctx context.Context, req *pb.YearAmountTransactionRequest) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	data, err := h.service.FindYearlyAmountSuccess(ctx, int(req.GetYear()))
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
		Message: "Successfully fetched yearly amount success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyAmountFailed(ctx context.Context, req *pb.MonthAmountTransactionRequest) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	request := &requests.MonthAmountTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyAmountFailed(ctx, request)
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
		Message: "Successfully fetched monthly amount failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyAmountFailed(ctx context.Context, req *pb.YearAmountTransactionRequest) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	data, err := h.service.FindYearlyAmountFailed(ctx, int(req.GetYear()))
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
		Message: "Successfully fetched yearly amount failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *pb.MonthMethodTransactionRequest) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodSuccess(ctx, request)
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
		Message: "Successfully fetched monthly transaction method success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyTransactionMethodSuccess(ctx context.Context, req *pb.YearMethodTransactionRequest) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.service.FindYearlyMethodSuccess(ctx, int(req.GetYear()))
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
		Message: "Successfully fetched yearly transaction method success stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetMonthlyTransactionMethodFailed(ctx context.Context, req *pb.MonthMethodTransactionRequest) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	request := &requests.MonthMethodTransaction{
		Year:  int(req.GetYear()),
		Month: int(req.GetMonth()),
	}

	data, err := h.service.FindMonthlyMethodFailed(ctx, request)
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
		Message: "Successfully fetched monthly transaction method failed stats",
		Data:    stats,
	}, nil
}

func (h *transactionStatsHandler) GetYearlyTransactionMethodFailed(ctx context.Context, req *pb.YearMethodTransactionRequest) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	data, err := h.service.FindYearlyMethodFailed(ctx, int(req.GetYear()))
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
		Message: "Successfully fetched yearly transaction method failed stats",
		Data:    stats,
	}, nil
}
