package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-order/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"

	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
)

type orderStatsHandler struct {
	pborder.UnimplementedOrderStatsServiceServer
	orderStats service.OrderStatsService
	logger     logger.LoggerInterface
}

func NewOrderStatsHandler(
	orderStats service.OrderStatsService,
	logger logger.LoggerInterface,
) OrderStatsHandler {
	return &orderStatsHandler{
		orderStats: orderStats,
		logger:     logger,
	}
}

func (s *orderStatsHandler) FindMonthlyTotalRevenue(ctx context.Context, req *pborder.FindYearMonthTotalRevenue) (*pborder.ApiResponseOrderMonthlyTotalRevenue, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthTotalRevenue{
		Year:  year,
		Month: month,
	}

	methods, err := s.orderStats.FindMonthlyTotalRevenue(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pborder.OrderMonthlyTotalRevenueResponse
	for _, method := range methods {
		data = append(data, &pborder.OrderMonthlyTotalRevenueResponse{
			Year:         method.Year,
			Month:        method.Month,
			TotalRevenue: int32(method.TotalRevenue),
		})
	}

	return &pborder.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    data,
	}, nil
}

func (s *orderStatsHandler) FindYearlyTotalRevenue(ctx context.Context, req *pborder.FindYearTotalRevenue) (*pborder.ApiResponseOrderYearlyTotalRevenue, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	methods, err := s.orderStats.FindYearlyTotalRevenue(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pborder.OrderYearlyTotalRevenueResponse
	for _, method := range methods {
		data = append(data, &pborder.OrderYearlyTotalRevenueResponse{
			Year:         method.Year,
			TotalRevenue: int32(method.TotalRevenue),
		})
	}

	return &pborder.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (s *orderStatsHandler) FindMonthlyRevenue(ctx context.Context, request *pborder.FindYearOrder) (*pborder.ApiResponseOrderMonthly, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindMonthlyOrder(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pborder.OrderMonthlyResponse
	for _, item := range res {
		data = append(data, &pborder.OrderMonthlyResponse{
			Month:          item.Month,
			OrderCount:     int32(item.OrderCount),
			TotalRevenue:   int32(item.TotalRevenue),
			TotalItemsSold: int32(item.TotalItemsSold),
		})
	}

	return &pborder.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue data retrieved",
		Data:    data,
	}, nil
}

func (s *orderStatsHandler) FindYearlyRevenue(ctx context.Context, request *pborder.FindYearOrder) (*pborder.ApiResponseOrderYearly, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindYearlyOrder(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var data []*pborder.OrderYearlyResponse
	for _, item := range res {
		data = append(data, &pborder.OrderYearlyResponse{
			Year:               item.Year,
			OrderCount:         int32(item.OrderCount),
			TotalRevenue:       int32(item.TotalRevenue),
			TotalItemsSold:     int32(item.TotalItemsSold),
			UniqueProductsSold: int32(item.UniqueProductsSold),
		})
	}

	return &pborder.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue data retrieved",
		Data:    data,
	}, nil
}

func (s *orderStatsHandler) mustEmbedUnimplementedOrderStatsServiceServer() {}
